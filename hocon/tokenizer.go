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

	arrayEndToken   = "]"
	arrayStartToken = "["

	endOfObjectToken = "}"
	objectStartToken = "{"

	commaToken   = ","
	dotToken     = "."
	newLineToken = `\n`

	plusAssignmentToken = "+="

	startOfQuotedTextToken       = `"`
	endOfQuotedTextToken         = `"`
	startOfQuotedKeyToken        = `"`
	endOfQuotedKeyToken          = `"`
	startOfTripleQuotedTextToken = `"""`
	endOfTripleQuotedTextToken   = `"""`

	escapeChar = `\`

	includeSpecial  = "include"
	optionalSpecial = '?'
)

var (
	assignmentTokens        = []string{"=", ":"}
	spaceOrTabTokens        = []string{" ", "\t"}
	startOfCommentTokens    = []string{"#", "//"}
	substitutionStartTokens = []string{"${", "${?"}
	//	HoconNotInUnquotedKey  = "$\"{}[]:=+,#`^?!@*&\\."
	unquotedKeyTokens = []string{"$", `"`, "{", "}", "[", "]", ":", "=",
		"+", ",", "#", "`", "^", "?", "!", "@", "*", "&", `\`, "."}

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
	whitespaceTokens = []string{
		" ", "\t", "\n", "\u000B", "\u000C",
		"\u000D", "\u00A0", "\u1680", "\u2000",
		"\u2001", "\u2002", "\u2003", "\u2004",
		"\u2005", "\u2006", "\u2007", "\u2008",
		"\u2009", "\u200A", "\u202F", "\u205F",
		"\u2060", "\u3000", "\uFEFF",
	}
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
	if p == nil {
		return false
	}
	return p.index >= len(p.text)
}

// Matches find any of the given patterns in tokenizer starting from the current peak,
// returns true when it is found, false - otherwise
func (p *Tokenizer) Matches(patterns ...string) bool {
	if p == nil {
		return false
	}

	for _, pattern := range patterns {
		if len(pattern)+p.index > len(p.text) {
			continue
		}

		selected := p.text[p.index : p.index+len(pattern)]

		if selected == pattern {
			return true
		}
	}

	return false
}

// MatchesMore find any of the given patterns in tokenizer starting from the current peak,
// returns true when it is found and followed by any other character, false - otherwise
func (p *Tokenizer) MatchesMore(patterns ...string) bool {
	if p == nil {
		return false
	}
	for _, pattern := range patterns {
		if len(pattern)+p.index >= len(p.text) { //
			continue
		}

		if p.text[p.index:p.index+len(pattern)] == pattern {
			return true
		}
	}
	return false
}

func (p *Tokenizer) Take(length int) string {
	if p == nil || p.index+length > len(p.text) {
		return ""
	}

	str := p.text[p.index : p.index+length]
	p.index += length
	return str
}

func (p *Tokenizer) Peek() byte {
	if p == nil || p.EOF() {
		return 0
	}

	return p.text[p.index]
}

func (p *Tokenizer) TakeOne() byte {
	if p == nil || p.EOF() {
		return 0
	}

	b := p.text[p.index]
	p.index += 1
	return b
}

func (p *HoconTokenizer) PullWhitespace() {
	for !p.EOF() && p.IsWhitespace() {
		p.TakeOne()
	}
}

type HoconTokenizer struct {
	*Tokenizer
}

func NewHoconTokenizer(text string) *HoconTokenizer {
	return &HoconTokenizer{NewTokenizer(text)}
}

func (p *HoconTokenizer) PullWhitespaceAndComments() {
	for {
		p.PullWhitespace()
		for p.IsStartOfComment() {
			p.PullComment()
		}

		if !p.IsWhitespace() {
			break
		}
	}
}

func (p *HoconTokenizer) PullRestOfLine() string {
	buf := bytes.NewBuffer(nil)

	for p.Tokenizer != nil && !p.EOF() {
		c := p.TakeOne()
		if c == '\n' {
			break
		}

		if c == '\r' {
			continue
		}
		if err := buf.WriteByte(c); err != nil {
			// Buffer.WriteByte never returns error
			panic(err)
		}
	}

	return strings.TrimSpace(buf.String())
}

func (p *HoconTokenizer) PullNext() (*Token, error) {
	var token *Token
	var err error

	p.PullWhitespaceAndComments()

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
		token = p.PullUnquotedKey()
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

	var msg string
	if p.Tokenizer == nil {
		msg = "unknown token"
	} else {
		msg = fmt.Sprintf("unknown token, offset: %d", p.index)
	}
	return nil, fmt.Errorf(msg)
}

func (p *HoconTokenizer) isStartOfQuotedKey() bool {
	return p.Matches(startOfQuotedKeyToken)
}

func (p *HoconTokenizer) PullArrayEnd() *Token {
	p.TakeOne()
	return NewToken(TokenTypeArrayEnd)
}

func (p *HoconTokenizer) IsArrayEnd() bool {
	return p.Matches(arrayEndToken)
}

func (p *HoconTokenizer) IsArrayStart() bool {
	return p.Matches(arrayStartToken)
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
	return p.Matches(commaToken)
}

func (p *HoconTokenizer) IsNewline() bool {
	return p.Matches(newLineToken)
}

func (p *HoconTokenizer) IsDot() bool {
	return p.Matches(dotToken)
}

func (p *HoconTokenizer) IsObjectStart() bool {
	return p.Matches(objectStartToken)
}

func (p *HoconTokenizer) IsEndOfObject() bool {
	return p.Matches(endOfObjectToken)
}

func (p *HoconTokenizer) IsAssignment() bool {
	return p.MatchesMore(assignmentTokens...)
}

func (p *HoconTokenizer) IsPlusAssignment() bool {
	return p.Matches(plusAssignmentToken)
}

func (p *HoconTokenizer) IsStartOfQuotedText() bool {
	return p.Matches(startOfQuotedTextToken)
}

func (p *HoconTokenizer) IsStartOfTripleQuotedText() bool {
	return p.Matches(startOfTripleQuotedTextToken)
}

func (p *HoconTokenizer) PullComment() *Token {
	p.PullRestOfLine()
	return NewToken(TokenTypeComment)
}

func (p *HoconTokenizer) PullUnquotedKey() *Token {
	buf := bytes.NewBuffer(nil)
	for !p.EOF() && p.IsUnquotedKey() {
		if err := buf.WriteByte(p.TakeOne()); err != nil {
			// Buffer.WriteByte never returns error
			panic(err)
		}
	}

	return NewTokenKey(strings.TrimSpace(buf.String()))
}

func (p *HoconTokenizer) IsUnquotedKey() bool {
	if p.Tokenizer == nil {
		return false
	}

	return !p.EOF() &&
		!p.IsStartOfComment() &&
		!p.Matches(unquotedKeyTokens...)
}

func (p *HoconTokenizer) IsUnquotedKeyStart() bool {
	if p.Tokenizer == nil {
		return false
	}

	return !p.EOF() &&
		!p.IsWhitespace() &&
		!p.IsStartOfComment() &&
		!p.Matches(unquotedKeyTokens...)
}

func (p *HoconTokenizer) IsWhitespace() bool {
	return p.Matches(whitespaceTokens...)
}

func (p *HoconTokenizer) IsWhitespaceOrComment() bool {
	return p.IsWhitespace() || p.IsStartOfComment()
}

func (p *HoconTokenizer) PullTripleQuotedText() (*Token, error) {
	if !p.IsStartOfTripleQuotedText() {
		return nil, fmt.Errorf("expected start of triple quoted text token, got %s", string(p.Peek()))
	}
	buf := bytes.NewBuffer(nil)
	p.Take(3)
	for !p.EOF() && !p.Matches(endOfTripleQuotedTextToken) {
		if err := buf.WriteByte(p.Peek()); err != nil {
			// Buffer.WriteByte cannot return error
			panic(err)
		}
		p.TakeOne()
	}
	p.Take(3)
	return NewTokenLiteralValue(buf.String()), nil
}

func (p *HoconTokenizer) PullQuotedText() (*Token, error) {
	if !p.IsStartOfQuotedText() {
		return nil, fmt.Errorf("expected start of quoted text token, got %s", string(p.Peek()))
	}
	buf := bytes.NewBuffer(nil)
	p.TakeOne()
	for !p.EOF() && !p.Matches(endOfQuotedTextToken) {
		if p.Matches(escapeChar) {
			sequence, err := p.pullEscapeSequence()
			if err != nil {
				return nil, err
			}

			if _, err := buf.WriteString(sequence); err != nil {
				// Buffer.WriteString cannot return error
				panic(err)
			}
		} else {
			if err := buf.WriteByte(p.Peek()); err != nil {
				// Buffer.WriteByte cannot return error
				panic(err)
			}
			p.TakeOne()
		}
	}
	p.TakeOne()
	return NewTokenLiteralValue(buf.String()), nil
}

func (p *HoconTokenizer) PullQuotedKey() (*Token, error) {
	if !p.isStartOfQuotedKey() {
		return nil, fmt.Errorf("expected start of quoted key token, got %s", string(p.Peek()))
	}
	buf := bytes.NewBuffer(nil)
	p.TakeOne()
	for !p.EOF() && !p.Matches(endOfQuotedKeyToken) {
		if p.Matches(escapeChar) {
			sequence, err := p.pullEscapeSequence()
			if err != nil {
				return nil, err
			}

			if _, err := buf.WriteString(sequence); err != nil {
				// Buffer.WriteString cannot return error
				panic(err)
			}
		} else {
			if err := buf.WriteByte(p.Peek()); err != nil {
				// Buffer.WriteString cannot return error
				panic(err)
			}
			p.TakeOne()
		}
	}
	p.TakeOne()
	return NewTokenKey(buf.String()), nil
}

func (p *HoconTokenizer) PullInclude() (*Token, error) {
	if !p.IsInclude() {
		return nil, fmt.Errorf("expected include token, got %s", string(p.Peek()))
	}
	p.Take(len(includeSpecial))
	p.PullWhitespaceAndComments()
	rest, err := p.PullQuotedText()
	if err != nil {
		return nil, err
	}

	unQuote := rest.value
	return NewTokenInclude(unQuote), nil
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
	return p.MatchesMore(startOfCommentTokens...)
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
		return p.pullUnquotedText(), nil
	}

	if p.IsArrayStart() {
		return p.PullArrayStart(), nil
	}

	if p.IsArrayEnd() {
		return p.PullArrayEnd(), nil
	}

	if p.IsSubstitutionStart() {
		return p.pullSubstitution(), nil
	}

	return nil, fmt.Errorf("expected value: Null literal, Array, Quoted Text, Unquoted Text, Triple quoted Text, Object or End of array")
}

func (p *HoconTokenizer) IsSubstitutionStart() bool {
	return p.MatchesMore(substitutionStartTokens...)
}

func (p *HoconTokenizer) IsInclude() bool {
	if p.Tokenizer == nil {
		return false
	}
	p.Push()
	defer func() {
		if err := p.Pop(); err != nil {
			panic(err)
		}
	}()
	if p.Matches(includeSpecial) {
		p.Take(len(includeSpecial))
		if p.IsWhitespaceOrComment() {
			p.PullWhitespaceAndComments()
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

func (p *HoconTokenizer) pullSubstitution() *Token {
	buf := bytes.NewBuffer(nil)
	p.Take(2)
	isOptional := false
	if p.Peek() == optionalSpecial {
		p.TakeOne()
		isOptional = true
	}

	for !p.EOF() && p.isUnquotedText() {
		if err := buf.WriteByte(p.TakeOne()); err != nil {
			// Buffer.WriteByte cannot return error
			panic(err)
		}
	}
	p.TakeOne()
	return NewTokenSubstitution(buf.String(), isOptional)
}

func (p *HoconTokenizer) IsSpaceOrTab() bool {
	return p.MatchesMore(spaceOrTabTokens...)
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

func (p *HoconTokenizer) PullSpaceOrTab() *Token {
	buf := bytes.NewBuffer(nil)
	for p.IsSpaceOrTab() {
		if err := buf.WriteByte(p.TakeOne()); err != nil {
			// Buffer.WriteByte cannot return error
			panic(err)
		}
	}
	return NewTokenLiteralValue(buf.String())
}

func (p *HoconTokenizer) pullUnquotedText() *Token {
	buf := bytes.NewBuffer(nil)
	for !p.EOF() && p.isUnquotedText() {
		if err := buf.WriteByte(p.TakeOne()); err != nil {
			// Buffer.WriteByte cannot return error
			panic(err)
		}
	}
	return NewTokenLiteralValue(buf.String())
}

func (p *HoconTokenizer) isUnquotedText() bool {
	if p.Tokenizer == nil {
		return false
	}
	return !p.EOF() &&
		!p.IsWhitespace() &&
		!p.IsStartOfComment() &&
		strings.IndexByte(HoconNotInUnquotedText, p.Peek()) == -1
}

func (p *HoconTokenizer) PullSimpleValue() (*Token, error) {
	if p.IsSpaceOrTab() {
		return p.PullSpaceOrTab(), nil
	}

	if p.isUnquotedText() {
		return p.pullUnquotedText(), nil
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
