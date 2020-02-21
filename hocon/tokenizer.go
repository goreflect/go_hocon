package hocon

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

const (
	HoconNotInUnquotedKey  = "$\"{}[]:=+,#`^?!@*&\\."
	HoconNotInUnquotedText = "$\"{}[]:=+,#`^?!@*&\\"
)

type Tokenizer struct {
	text       string
	index      int
	indexStack *Stack
}

func NewTokenizer(text string) *Tokenizer {
	return &Tokenizer{
		indexStack: NewStack(),
		text:       text,
	}
}

func (p *Tokenizer) Push() {
	p.indexStack.Push(p.index)
}

func (p *Tokenizer) Pop() error {
	index, err := p.indexStack.Pop()
	if err != nil {
		return err
	}

	p.index = index
	return nil
}

func (p *Tokenizer) EOF() bool {
	return p.index >= len(p.text)
}

func (p *Tokenizer) Matches(pattern string) bool {

	if len(pattern)+p.index > len(p.text) {
		return false
	}

	selected := p.text[p.index : p.index+len(pattern)]

	return selected == pattern
}

func (p *Tokenizer) MatchesMore(patterns []string) bool {
	for _, pattern := range patterns {
		if len(pattern)+p.index >= len(p.text) {
			continue
		}

		if p.text[p.index:p.index+len(pattern)] == pattern {
			return true
		}
	}
	return false
}

func (p *Tokenizer) Take(length int) string {
	if p.index+length > len(p.text) {
		return ""
	}

	str := p.text[p.index : p.index+length]
	p.index += length
	return str
}

func (p *Tokenizer) Peek() byte {
	if p.EOF() {
		return 0
	}

	return p.text[p.index]
}

func (p *Tokenizer) TakeOne() byte {
	if p.EOF() {
		return 0
	}

	b := p.text[p.index]
	p.index += 1
	return b
}

func (p *Tokenizer) PullWhitespace() {
	for !p.EOF() && isWhitespace(p.Peek()) {
		p.TakeOne()
	}
}

type HoconTokenizer struct {
	*Tokenizer
}

func NewHoconTokenizer(text string) *HoconTokenizer {
	return &HoconTokenizer{NewTokenizer(text)}
}

func (p *HoconTokenizer) PullWhitespaceAndComments() error {
	for {
		p.PullWhitespace()
		for p.IsStartOfComment() {
			if _, err := p.PullComment(); err != nil {
				return err
			}
		}

		if !p.IsWhitespace() {
			break
		}
	}
	return nil
}

func (p *HoconTokenizer) PullRestOfLine() (string, error) {
	buf := bytes.NewBuffer(nil)

	for !p.EOF() {
		c := p.TakeOne()
		if c == '\n' {
			break
		}

		if c == '\r' {
			continue
		}
		if err := buf.WriteByte(c); err != nil {
			return "", err
		}
	}

	return strings.TrimSpace(buf.String()), nil
}

func (p *HoconTokenizer) PullNext() (*Token, error) {
	var token *Token
	var err error

	if err := p.PullWhitespaceAndComments(); err != nil {
		return nil, err
	}
	if p.IsDot() {
		token = p.PullDot()
	} else if p.IsObjectStart() {
		token = p.PullStartOfObject()
	} else if p.IsEndOfObject() {
		token = p.PullEndOfObject()
	} else if p.IsAssignment() {
		token = p.PullAssignment()
	} else if p.IsPlusAssignment() {
		token = p.PullPlusAssignment()
	} else if p.IsInclude() {
		token, err = p.PullInclude()
		if err != nil {
			return nil, err
		}
	} else if p.isStartOfQuotedKey() {
		token, err = p.PullQuotedKey()
		if err != nil {
			return nil, err
		}
	} else if p.IsUnquotedKeyStart() {
		token, err = p.PullUnquotedKey()
		if err != nil {
			return nil, err
		}
	} else if p.IsArrayStart() {
		token = p.PullArrayStart()
	} else if p.IsArrayEnd() {
		token = p.PullArrayEnd()
	} else if p.EOF() {
		token = NewToken(TokenTypeEoF)
	}

	if token != nil {
		return token, nil
	}

	return nil, fmt.Errorf("unknown token, offset: %d", p.index)
}

func (p *HoconTokenizer) isStartOfQuotedKey() bool {
	return p.Matches("\"")
}

func (p *HoconTokenizer) PullArrayEnd() *Token {
	p.TakeOne()
	return NewToken(TokenTypeArrayEnd)
}

func (p *HoconTokenizer) IsArrayEnd() bool {
	return p.Matches("]")
}

func (p *HoconTokenizer) IsArrayStart() bool {
	return p.Matches("[")
}

func (p *HoconTokenizer) PullArrayStart() *Token {
	p.TakeOne()
	return NewToken(TokenTypeArrayStart)
}

func (p *HoconTokenizer) PullDot() *Token {
	p.TakeOne()
	return NewToken(TokenTypeDot)
}

func (p *HoconTokenizer) PullComma() *Token {
	p.TakeOne()
	return NewToken(TokenTypeComma)
}

func (p *HoconTokenizer) PullNewline() *Token {
	p.Take(2)
	return NewToken(TokenTypeNewline)
}

func (p *HoconTokenizer) PullStartOfObject() *Token {
	p.TakeOne()
	return NewToken(TokenTypeObjectStart)
}

func (p *HoconTokenizer) PullEndOfObject() *Token {
	p.TakeOne()
	return NewToken(TokenTypeObjectEnd)
}

func (p *HoconTokenizer) PullAssignment() *Token {
	p.TakeOne()
	return NewToken(TokenTypeAssign)
}

func (p *HoconTokenizer) PullPlusAssignment() *Token {
	p.Take(2)
	return NewToken(TokenTypePlusAssign)
}

func (p *HoconTokenizer) IsComma() bool {
	return p.Matches(",")
}

func (p *HoconTokenizer) IsNewline() bool {
	return p.Matches(`\n`)
}

func (p *HoconTokenizer) IsDot() bool {
	return p.Matches(".")
}

func (p *HoconTokenizer) IsObjectStart() bool {
	return p.Matches("{")
}

func (p *HoconTokenizer) IsEndOfObject() bool {
	return p.Matches("}")
}

func (p *HoconTokenizer) IsAssignment() bool {
	return p.MatchesMore([]string{"=", ":"})
}

func (p *HoconTokenizer) IsPlusAssignment() bool {
	return p.Matches("+=")
}

func (p *HoconTokenizer) IsStartOfQuotedText() bool {
	return p.Matches("\"")
}

func (p *HoconTokenizer) IsStartOfTripleQuotedText() bool {
	return p.Matches("\"\"\"")
}

func (p *HoconTokenizer) PullComment() (*Token, error) {
	if _, err := p.PullRestOfLine(); err != nil {
		return nil, err
	}
	return NewToken(TokenTypeComment), nil
}

func (p *HoconTokenizer) PullUnquotedKey() (*Token, error) {
	buf := bytes.NewBuffer(nil)
	for !p.EOF() && p.IsUnquotedKey() {
		if err := buf.WriteByte(p.TakeOne()); err != nil {
			return nil, err
		}
	}

	return DefaultToken.Key(strings.TrimSpace(buf.String())), nil
}

func (p *HoconTokenizer) IsUnquotedKey() bool {
	return !p.EOF() && !p.IsStartOfComment() && (strings.IndexByte(HoconNotInUnquotedKey, p.Peek()) == -1)
}

func (p *HoconTokenizer) IsUnquotedKeyStart() bool {
	return !p.EOF() && !p.IsWhitespace() && !p.IsStartOfComment() && (strings.IndexByte(HoconNotInUnquotedKey, p.Peek()) == -1)
}

func (p *HoconTokenizer) IsWhitespace() bool {
	return isWhitespace(p.Peek())
}

func (p *HoconTokenizer) IsWhitespaceOrComment() bool {
	return p.IsWhitespace() || p.IsStartOfComment()
}

func (p *HoconTokenizer) PullTripleQuotedText() (*Token, error) {
	buf := bytes.NewBuffer(nil)
	p.Take(3)
	for !p.EOF() && !p.Matches("\"\"\"") {
		if err := buf.WriteByte(p.Peek()); err != nil {
			return nil, err
		}
		p.TakeOne()
	}
	p.Take(3)
	return DefaultToken.LiteralValue(buf.String()), nil
}

func (p *HoconTokenizer) PullQuotedText() (*Token, error) {
	buf := bytes.NewBuffer(nil)
	p.TakeOne()
	for !p.EOF() && !p.Matches("\"") {
		if p.Matches("\\") {
			sequence, err := p.pullEscapeSequence()
			if err != nil {
				return nil, err
			}

			if _, err := buf.WriteString(sequence); err != nil {
				return nil, err
			}
		} else {
			if err := buf.WriteByte(p.Peek()); err != nil {
				return nil, err
			}
			p.TakeOne()
		}
	}
	p.TakeOne()
	return DefaultToken.LiteralValue(buf.String()), nil
}

func (p *HoconTokenizer) PullQuotedKey() (*Token, error) {
	buf := bytes.NewBuffer(nil)
	p.TakeOne()
	for !p.EOF() && !p.Matches("\"") {
		if p.Matches("\\") {
			sequence, err := p.pullEscapeSequence()
			if err != nil {
				return nil, err
			}

			if _, err := buf.WriteString(sequence); err != nil {
				return nil, err
			}
		} else {
			if err := buf.WriteByte(p.Peek()); err != nil {
				return nil, err
			}
			p.TakeOne()
		}
	}
	p.TakeOne()
	return DefaultToken.Key(buf.String()), nil
}

func (p *HoconTokenizer) PullInclude() (*Token, error) {
	p.Take(len("include"))
	if err := p.PullWhitespaceAndComments(); err != nil {
		return nil, err
	}
	rest, err := p.PullQuotedText()
	if err != nil {
		return nil, err
	}

	unQuote := rest.value
	return DefaultToken.Include(unQuote), nil
}

func (p *HoconTokenizer) pullEscapeSequence() (string, error) {
	p.TakeOne()
	escaped := p.TakeOne()
	switch escaped {
	case '"':
		return "\"", nil
	case '\\':
		return "\\", nil
	case '/':
		return "/", nil
	case 'b':
		return "\b", nil
	case 'f':
		return "\f", nil
	case 'n':
		return "\n", nil
	case 'r':
		return "\r", nil
	case 't':
		return "\t", nil
	case 'u':
		utf8Code := "\\u" + strings.ToLower(p.Take(4))
		utf8Str := ""
		if _, err := fmt.Sscanf(utf8Code, "%s", &utf8Str); err != nil {
			return "", err
		}
		return utf8Str, nil
	default:
		return "", fmt.Errorf("unknown escape code: %v", escaped)
	}
}

func (p *HoconTokenizer) IsStartOfComment() bool {
	return p.MatchesMore([]string{"#", "//"})
}

func (p *HoconTokenizer) PullValue() (*Token, error) {
	if p.IsObjectStart() {
		return p.PullStartOfObject(), nil
	}

	if p.IsStartOfTripleQuotedText() {
		return p.PullTripleQuotedText()
	}

	if p.IsStartOfQuotedText() {
		return p.PullQuotedText()
	}

	if p.isUnquotedText() {
		return p.pullUnquotedText()
	}

	if p.IsArrayStart() {
		return p.PullArrayStart(), nil
	}

	if p.IsArrayEnd() {
		return p.PullArrayEnd(), nil
	}

	if p.IsSubstitutionStart() {
		return p.pullSubstitution()
	}

	return nil, fmt.Errorf("expected value: Null literal, Array, Quoted Text, Unquoted Text, Triple quoted Text, Object or End of array")
}

func (p *HoconTokenizer) IsSubstitutionStart() bool {
	return p.MatchesMore([]string{"${", "${?"})
}

func (p *HoconTokenizer) IsInclude() bool {
	p.Push()
	defer func() {
		if err := p.Pop(); err != nil {
			panic(err)
		}
	}()
	if p.Matches("include") {
		p.Take(len("include"))
		if p.IsWhitespaceOrComment() {
			if err := p.PullWhitespaceAndComments(); err != nil {
				return false
			}
			if p.IsStartOfQuotedText() {
				if _, err := p.PullQuotedText(); err != nil {
					return false
				}
				return true
			}
		}
	}

	return false
}

func (p *HoconTokenizer) pullSubstitution() (*Token, error) {
	buf := bytes.NewBuffer(nil)
	p.Take(2)
	isOptional := false
	if p.Peek() == '?' {
		p.TakeOne()
		isOptional = true
	}

	for !p.EOF() && p.isUnquotedText() {
		if err := buf.WriteByte(p.TakeOne()); err != nil {
			return nil, err
		}
	}
	p.TakeOne()
	return DefaultToken.Substitution(buf.String(), isOptional), nil
}

func (p *HoconTokenizer) IsSpaceOrTab() bool {
	return p.MatchesMore([]string{" ", "\t"})
}

func (p *HoconTokenizer) IsStartSimpleValue() bool {
	if p.IsSpaceOrTab() {
		return true
	}

	if p.isUnquotedText() {
		return true
	}

	return false
}

func (p *HoconTokenizer) PullSpaceOrTab() (*Token, error) {
	buf := bytes.NewBuffer(nil)
	for p.IsSpaceOrTab() {
		if err := buf.WriteByte(p.TakeOne()); err != nil {
			return nil, err
		}
	}
	return DefaultToken.LiteralValue(buf.String()), nil
}

func (p *HoconTokenizer) pullUnquotedText() (*Token, error) {
	buf := bytes.NewBuffer(nil)
	for !p.EOF() && p.isUnquotedText() {
		if err := buf.WriteByte(p.TakeOne()); err != nil {
			return nil, err
		}
	}
	return DefaultToken.LiteralValue(buf.String()), nil
}

func (p *HoconTokenizer) isUnquotedText() bool {
	return !p.EOF() && !p.IsWhitespace() && !p.IsStartOfComment() && strings.IndexByte(HoconNotInUnquotedText, p.Peek()) == -1
}

func (p *HoconTokenizer) PullSimpleValue() (*Token, error) {
	if p.IsSpaceOrTab() {
		return p.PullSpaceOrTab()
	}

	if p.isUnquotedText() {
		return p.pullUnquotedText()
	}
	return nil, errors.New("no simple value found")
}

func (p *HoconTokenizer) isValue() bool {

	if p.IsArrayStart() ||
		p.IsObjectStart() ||
		p.IsStartOfTripleQuotedText() ||
		p.IsSubstitutionStart() ||
		p.IsStartOfQuotedText() ||
		p.isUnquotedText() {
		return true
	}
	return false
}

/*
SPACE (\u0020)
NO-BREAK SPACE (\u00A0)
OGHAM SPACE MARK (\u1680)
EN QUAD (\u2000)
EM QUAD (\u2001)
EN SPACE (\u2002)
EM SPACE (\u2003)
THREE-PER-EM SPACE (\u2004)
FOUR-PER-EM SPACE (\u2005)
SIX-PER-EM SPACE (\u2006)
FIGURE SPACE (\u2007)
PUNCTUATION SPACE (\u2008)
THIN SPACE (\u2009)
HAIR SPACE (\u200A)
NARROW NO-BREAK SPACE (\u202F)
MEDIUM MATHEMATICAL SPACE (\u205F)
and IDEOGRAPHIC SPACE (\u3000)
Byte Order Mark (\uFEFF)
*/
func isWhitespace(c byte) bool {
	str := string(c)

	switch str {
	case " ", "\t", "\n", "\u000B", "\u000C",
		"\u000D", "\u00A0", "\u1680", "\u2000",
		"\u2001", "\u2002", "\u2003", "\u2004",
		"\u2005", "\u2006", "\u2007", "\u2008",
		"\u2009", "\u200A", "\u202F", "\u205F",
		"\u2060", "\u3000", "\uFEFF":
		return true
	}
	return false
}
