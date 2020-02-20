package hocon

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type IncludeCallback func(filename string) (*HoconRoot, error)

type Parser struct {
	reader   *HoconTokenizer
	root     *HoconValue
	callback IncludeCallback

	substitutions []*HoconSubstitution
}

func Parse(text string, callback IncludeCallback) (*HoconRoot, error) {
	return new(Parser).parseText(text, callback)
}

func (p *Parser) parseText(text string, callback IncludeCallback) (*HoconRoot, error) {
	p.callback = callback
	p.root = NewHoconValue()
	p.reader = NewHoconTokenizer(text)
	if err := p.reader.PullWhitespaceAndComments(); err != nil {
		return nil, err
	}

	if err := p.parseObject(p.root, true, ""); err != nil {
		return nil, err
	}

	root := NewHoconRoot(p.root)

	cRoot := root.Value()

	for _, sub := range p.substitutions {
		res, err := getNode(cRoot, sub.Path)
		if err != nil {
			return nil, err
		}

		if res == nil {
			envVal, exist := os.LookupEnv(sub.OriginalPath)
			if !exist {
				if !sub.IsOptional {
					return nil, fmt.Errorf("unresolved substitution: %s", sub.Path)
				}
			} else {
				hv := NewHoconValue()
				hv.AppendValue(NewHoconLiteral(envVal))
				sub.ResolvedValue = hv
			}
		} else {
			sub.ResolvedValue = res
		}
	}

	return NewHoconRoot(p.root, p.substitutions...), nil
}

func (p *Parser) parseObject(owner *HoconValue, root bool, currentPath string) error {
	if !owner.IsObject() {
		owner.NewValue(NewHoconObject())
	}

	if owner.IsObject() {
		rootObj := owner
		for rootObj.oldValue != nil {
			oldObj, err := rootObj.oldValue.GetObject()
			if err != nil {
				return err
			}

			obj, err := rootObj.GetObject()
			if err != nil {
				return err
			}

			if oldObj == nil || obj == nil {
				break
			}
			if err := obj.Merge(oldObj); err != nil {
				return err
			}
			rootObj = rootObj.oldValue
		}
	}

	currentObject, err := owner.GetObject()
	if err != nil {
		return err
	}

	for !p.reader.EOF() {
		t, err := p.reader.PullNext()
		if err != nil {
			return err
		}

		switch t.tokenType {
		case TokenTypeInclude:
			included, err := p.callback(t.value)
			if err != nil {
				return err
			}

			substitutions := included.substitutions
			for _, substitution := range substitutions {
				substitution.Path = currentPath + "." + substitution.Path
			}
			p.substitutions = append(p.substitutions, substitutions...)
			otherObj, err := included.value.GetObject()
			if err != nil {
				return err
			}

			objectV, err := owner.GetObject()
			if err != nil {
				return err
			}

			if err := objectV.Merge(otherObj); err != nil {
				return err
			}
		case TokenTypeEoF:
		case TokenTypeKey:
			value := currentObject.GetOrCreateKey(t.value)
			nextPath := t.value
			if len(currentPath) > 0 {
				nextPath = currentPath + "." + t.value
			}
			if err := p.parseKeyContent(value, nextPath); err != nil {
				return err
			}
			if !root {
				return nil
			}
		case TokenTypeObjectEnd:
			return nil
		}
	}
	return nil
}

func (p *Parser) parseKeyContent(value *HoconValue, currentPath string) error {
	for !p.reader.EOF() {
		t, err := p.reader.PullNext()
		if err != nil {
			return err
		}

		switch t.tokenType {
		case TokenTypeDot:
			return p.parseObject(value, false, currentPath)
		case TokenTypeAssign:
			{
				if !value.IsObject() {
					value.Clear()
				}
			}
			return p.ParseValue(value, false, currentPath)
		case TokenTypePlusAssign:
			{
				if !value.IsObject() {
					value.Clear()
				}
			}
			return p.ParseValue(value, true, currentPath)
		case TokenTypeObjectStart:
			return p.parseObject(value, true, currentPath)
		}
	}
	return nil
}

func (p *Parser) ParseValue(owner *HoconValue, isEqualPlus bool, currentPath string) error {
	if p.reader.EOF() {
		return errors.New("end of file reached while trying to read a value")
	}

	if err := p.reader.PullWhitespaceAndComments(); err != nil {
		return err
	}

	for p.reader.isValue() {
		t, err := p.reader.PullValue()
		if err != nil {
			return err
		}

		if isEqualPlus {
			sub := p.ParseSubstitution(currentPath, false)
			p.substitutions = append(p.substitutions, sub)
			owner.AppendValue(sub)
		}

		switch t.tokenType {
		case TokenTypeEoF:
		case TokenTypeLiteralValue:
			if owner.IsObject() {
				owner.Clear()
			}
			lit := NewHoconLiteral(t.value)
			owner.AppendValue(lit)
		case TokenTypeObjectStart:
			if err := p.parseObject(owner, true, currentPath); err != nil {
				return err
			}
		case TokenTypeArrayStart:
			arr, err := p.ParseArray(currentPath)
			if err != nil {
				return err
			}

			owner.AppendValue(&arr)
		case TokenTypeSubstitute:
			sub := p.ParseSubstitution(t.value, t.isOptional)
			p.substitutions = append(p.substitutions, sub)
			owner.AppendValue(sub)
		}

		if p.reader.IsSpaceOrTab() {
			if err := p.ParseTrailingWhitespace(owner); err != nil {
				return err
			}
		}
	}
	p.ignoreComma()
	p.ignoreNewline()
	return nil
}

func (p *Parser) ParseTrailingWhitespace(owner *HoconValue) error {
	ws, err := p.reader.PullSpaceOrTab()
	if err != nil {
		return err
	}

	if len(ws.value) > 0 {
		wsList := NewHoconLiteral(ws.value)
		owner.AppendValue(wsList)
	}
	return nil
}

func (p *Parser) ParseSubstitution(value string, isOptional bool) *HoconSubstitution {
	return NewHoconSubstitution(value, isOptional)
}

func (p *Parser) ParseArray(currentPath string) (HoconArray, error) {
	arr := NewHoconArray()
	for !p.reader.EOF() && !p.reader.IsArrayEnd() {
		v := NewHoconValue()
		if err := p.ParseValue(v, false, currentPath); err != nil {
			return HoconArray{}, err
		}
		arr.values = append(arr.values, v)
		if err := p.reader.PullWhitespaceAndComments(); err != nil {
			return HoconArray{}, err
		}
	}
	p.reader.PullArrayEnd()
	return *arr, nil
}

func (p *Parser) ignoreComma() {
	if p.reader.IsComma() {
		p.reader.PullComma()
	}
}

func (p *Parser) ignoreNewline() {
	if p.reader.IsNewline() {
		p.reader.PullNewline()
	}
}

func getNode(root *HoconValue, path string) (*HoconValue, error) {
	elements := splitDottedPathHonouringQuotes(path)
	currentNode := root

	if currentNode == nil {
		return nil, errors.New("current node should not be null")
	}

	for _, key := range elements {
		var err error
		currentNode, err = currentNode.GetChildObject(key)
		if err != nil {
			return nil, err
		}

		if currentNode == nil {
			return nil, nil
		}
	}
	return currentNode, nil
}

func splitDottedPathHonouringQuotes(path string) []string {
	tmp1 := strings.Split(path, "\"")
	var values []string
	for i := 0; i < len(tmp1); i++ {
		tmp2 := strings.Split(tmp1[i], ".")
		for j := 0; j < len(tmp2); j++ {
			if len(tmp2[j]) > 0 {
				values = append(values, tmp2[j])
			}
		}
	}
	return values
}
