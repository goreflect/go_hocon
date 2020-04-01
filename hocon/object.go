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
	if p.items == nil || len(p.items) == 0 {
		return nil, fmt.Errorf("empty HoconObject cannot be unwrapped")
	}

	dict := map[string]interface{}{}

	for _, k := range p.keys {
		v := p.items[k]

		if v.IsObject() {
			obj, err1 := v.GetObject()
			// must not return error after checking v.IsObject()
			if err1 != nil {
				panic(err1)
			}

			unwrapped, err2 := obj.Unwrapped()
			// must not return error after checking v.IsObject()
			if err2 != nil {
				panic(err1)
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
	value := p.items[key]
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
	return p.ToString(0)
}

func (p *HoconObject) ToString(indent int) string {
	tmp := strings.Repeat(" ", indent*2)
	buf := bytes.NewBuffer(nil)
	for _, k := range p.keys {
		key := quoteStringIfNeeded(k)
		v := p.items[key]

		str := v.ToString(indent)
		buf.WriteString(fmt.Sprintf("%s%s : %s%s", tmp, key, str, newLine))
	}
	return buf.String()
}

func (p *HoconObject) Merge(other *HoconObject) {
	if other == nil {
		return
	}

	thisValues := p.items
	otherItems := other.items

	otherKeys := other.keys

	for _, otherKey := range otherKeys {
		otherValue := otherItems[otherKey]

		if thisValue, exist := thisValues[otherKey]; exist {
			isThisObject := thisValue.IsObject()
			isOtherObject := otherValue.IsObject()

			if isThisObject && isOtherObject {
				thisValueObject, err := thisValue.GetObject()
				// must not return error after checking thisValue.IsObject()
				if err != nil {
					panic(err)
				}

				otherObjectValue, err := otherValue.GetObject()
				// must not return error after checking otherValue.IsObject()
				if err != nil {
					panic(err)
				}

				thisValueObject.Merge(otherObjectValue)
			}
		} else {
			p.items[otherKey] = otherValue
			p.keys = append(p.keys, otherKey)
		}
	}
	return
}

func (p *HoconObject) MergeImmutable(other *HoconObject) *HoconObject {
	thisValues := p.items
	thisKeys := p.keys
	newObject := HoconObject{items: thisValues, keys: thisKeys}

	if other == nil {
		return &newObject
	}

	newObject.Merge(other)

	return &newObject
}
