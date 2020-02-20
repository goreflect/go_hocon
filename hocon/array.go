package hocon

import (
	"errors"
	"strings"
)

type HoconArray struct {
	values []*HoconValue
}

func NewHoconArray() *HoconArray {
	return &HoconArray{}
}

func (p *HoconArray) IsString() bool {
	return false
}

func (p *HoconArray) GetString() (string, error) {
	return "", errors.New("this element is an array and not a string")
}

func (p *HoconArray) IsArray() bool {
	return true
}

func (p *HoconArray) GetArray() ([]*HoconValue, error) {
	return p.values, nil
}

func (p *HoconArray) String() string {
	var sstr []string
	for _, v := range p.values {
		sstr = append(sstr, v.String())
	}
	return "[" + strings.Join(sstr, ",") + "]"
}
