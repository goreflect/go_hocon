package hocon

import (
	"fmt"
)

type HoconSubstitution struct {
	Path          string
	ResolvedValue *HoconValue
	IsOptional    bool
	OriginalPath  string
}

func NewHoconSubstitution(path string, isOptional bool) *HoconSubstitution {
	return &HoconSubstitution{Path: path, OriginalPath: path, IsOptional: isOptional}
}

func (p *HoconSubstitution) IsString() bool {
	if p.ResolvedValue == nil {
		return false
	}
	if err := p.checkCycleRef(); err != nil {
		return false
	}
	return p.ResolvedValue.IsString()
}

func (p *HoconSubstitution) GetString() (string, error) {
	if p.ResolvedValue == nil {
		return "", nil
	}
	if err := p.checkCycleRef(); err != nil {
		return "", err
	}
	return p.ResolvedValue.GetString()
}

func (p *HoconSubstitution) IsArray() bool {
	if p.ResolvedValue == nil {
		return false
	}
	if err := p.checkCycleRef(); err != nil {
		return false
	}
	return p.ResolvedValue.IsArray()
}
func (p *HoconSubstitution) GetArray() ([]*HoconValue, error) {
	if p.ResolvedValue == nil {
		return nil, nil
	}
	return p.ResolvedValue.GetArray()
}

func (p *HoconSubstitution) IsObject() (bool, error) {
	if p.ResolvedValue == nil {
		return false, nil
	}
	if err := p.checkCycleRef(); err != nil {
		return false, err
	}
	return p.ResolvedValue.IsObject()
}

func (p *HoconSubstitution) GetObject() (*HoconObject, error) {
	if p.ResolvedValue == nil {
		return nil, nil
	}
	if err := p.checkCycleRef(); err != nil {
		return nil, err
	}
	return p.ResolvedValue.GetObject()
}

func (p *HoconSubstitution) checkCycleRef() error {
	if p.hasCycleRef(map[HoconElement]int{}, 1) {
		return fmt.Errorf("cycle reference in path of %s", p.Path)
	}
	return nil
}

func (p *HoconSubstitution) hasCycleRef(dup map[HoconElement]int, level int) bool {
	if p.ResolvedValue == nil {
		return false
	}

	if lvl, exist := dup[p.ResolvedValue]; exist {
		if lvl != level {
			return true
		}
	}
	dup[p.ResolvedValue] = level

	for _, subV := range p.ResolvedValue.values {
		if sub, ok := subV.(*HoconSubstitution); ok {
			if sub.ResolvedValue != nil {
				return sub.hasCycleRef(dup, level+1)
			}
		}
	}

	return false
}
