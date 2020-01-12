package hocon

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

type HoconObject struct {
	items map[string]*HoconValue
	keys  []string
}

func NewHoconObject() *HoconObject {
	return &HoconObject{
		items: make(map[string]*HoconValue),
	}
}

// makeHoconObject creates object with text values for test purposes
func makeHoconObject(keys []string, values []string) *HoconObject {
	items := make(map[string]*HoconValue)
	for k, v := range keys {
		items[v] = &HoconValue{values: []HoconElement{NewHoconLiteral(values[k])}}
	}

	return &HoconObject{
		keys:  keys,
		items: items,
	}
}

func (p *HoconObject) GetString() (string, error) {
	return "", errors.New("this element is an object and not a string")
}

func (p *HoconObject) IsArray() bool {
	return false
}

func (p *HoconObject) GetArray() ([]*HoconValue, error) {
	return nil, errors.New("this element is an object and not an array")
}

func (p *HoconObject) GetKeys() []string {
	return p.keys
}

func (p *HoconObject) Unwrapped() (map[string]interface{}, error) {
	if len(p.items) == 0 {
		return nil, nil
	}

	dict := map[string]interface{}{}

	for _, k := range p.keys {
		v := p.items[k]

		obj, err := v.GetObject()
		if err != nil {
			return nil, err
		}
		if obj != nil {
			unwrapped, err := obj.Unwrapped()
			if err != nil {
				return nil, err
			}
			dict[k] = unwrapped
		} else {
			dict[k] = v
		}
	}

	return dict, nil
}

func (p *HoconObject) Items() map[string]*HoconValue {
	return p.items
}

func (p *HoconObject) GetKey(key string) *HoconValue {
	value, _ := p.items[key]
	return value
}

func (p *HoconObject) GetOrCreateKey(key string) *HoconValue {
	if value, exist := p.items[key]; exist {
		child := NewHoconValue()
		child.oldValue = value
		p.items[key] = child
		return child
	}

	child := NewHoconValue()
	if p.items == nil {
		p.items = map[string]*HoconValue{}
	}
	p.items[key] = child
	p.keys = append(p.keys, key)
	return child
}

func (p *HoconObject) IsString() bool {
	return false
}

func (p *HoconObject) String() string {
	str, err := p.ToString(0)
	if err != nil {
		return fmt.Sprintf("cannot get string: %s", err.Error())
	}
	return str
}

func (p *HoconObject) ToString(indent int) (string, error) {
	tmp := strings.Repeat(" ", indent*2)
	buf := bytes.NewBuffer(nil)
	for _, k := range p.keys {
		key := p.quoteIfNeeded(k)
		v := p.items[key]

		str, err := v.ToString(indent)
		if err != nil {
			return "", err
		}
		buf.WriteString(fmt.Sprintf("%s%s : %s\r\n", tmp, key, str))
	}
	return buf.String(), nil
}

func (p *HoconObject) Merge(other *HoconObject) error {
	thisValues := p.items
	otherItems := other.items

	otherKeys := other.keys

	for _, otherKey := range otherKeys {
		otherValue := otherItems[otherKey]

		if thisValue, exist := thisValues[otherKey]; exist {
			if thisValue.IsObject() && otherValue.IsObject() {
				thisValueObject, err := thisValue.GetObject()
				if err != nil {
					return err
				}
				otherObjectValue, err := otherValue.GetObject()
				if err != nil {
					return err
				}
				if err := thisValueObject.Merge(otherObjectValue); err != nil {
					return err
				}
			}
		} else {
			p.items[otherKey] = otherValue
			p.keys = append(p.keys, otherKey)
		}
	}
	return nil
}

func (p *HoconObject) MergeImmutable(other *HoconObject) (*HoconObject, error) {
	thisValues := map[string]*HoconValue{}
	otherKeys := other.keys

	var thisKeys []string

	otherItems := other.items

	for _, otherKey := range otherKeys {
		otherValue := otherItems[otherKey]

		if thisValue, exist := thisValues[otherKey]; exist {

			if thisValue.IsObject() && otherValue.IsObject() {
				thisValueObject, err := thisValue.GetObject()
				if err != nil {
					return nil, err
				}

				otherValueObject, err := otherValue.GetObject()
				if err != nil {
					return nil, err
				}

				mergedObject, err := thisValueObject.MergeImmutable(otherValueObject)
				if err != nil {
					return nil, err
				}
				mergedValue := NewHoconValue()

				mergedValue.AppendValue(mergedObject)
				thisValues[otherKey] = mergedValue
			}
		} else {
			thisValues[otherKey] = &HoconValue{values: otherValue.values}
			thisKeys = append(thisKeys, otherKey)
		}
	}

	return &HoconObject{items: thisValues, keys: thisKeys}, nil
}

func (p *HoconObject) quoteIfNeeded(text string) string {
	if strings.IndexByte(text, ' ') >= 0 ||
		strings.IndexByte(text, '\t') >= 0 {
		return "\"" + text + "\""
	}
	return text
}
