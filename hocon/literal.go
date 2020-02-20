package hocon

import "errors"

type HoconLiteral struct {
	value string
}

func NewHoconLiteral(value string) *HoconLiteral {
	return &HoconLiteral{value: value}
}

func (p *HoconLiteral) IsString() bool {
	return true
}

func (p *HoconLiteral) GetString() (string, error) {
	return p.value, nil
}

func (p *HoconLiteral) IsArray() bool {
	return false
}

func (p *HoconLiteral) GetArray() ([]*HoconValue, error) {
	return nil, errors.New("this element is a string literal and not an array")
}

func (p *HoconLiteral) String() string {
	return p.value
}
