package configuration

import (
	"errors"
	"math/big"
	"strings"
	"time"

	"github.com/goreflect/go_hocon/hocon"
)

type Config struct {
	root          *hocon.HoconValue
	substitutions []*hocon.HoconSubstitution
	fallback      *Config
}

func NewConfigFromRoot(root *hocon.HoconRoot) (*Config, error) {
	if root.Value() == nil {
		return nil, errors.New("the root value cannot be null")
	}

	return &Config{
		root:          root.Value(),
		substitutions: root.Substitutions(),
	}, nil
}

func NewConfigFromConfig(source, fallback *Config) (*Config, error) {
	if source == nil {
		return nil, errors.New("the source configuration cannot be null")
	}

	return &Config{
		root:     source.root,
		fallback: fallback,
	}, nil
}

func (p *Config) IsEmpty() bool {
	return p == nil || p.root == nil || p.root.IsEmpty()
}

func (p *Config) Root() *hocon.HoconValue {
	return p.root
}

func (p *Config) Copy(fallback ...*Config) *Config {

	var fb *Config

	if p.fallback != nil {
		fb = p.fallback.Copy()
	} else {
		if len(fallback) > 0 {
			fb = fallback[0]
		}
	}
	return &Config{
		fallback:      fb,
		root:          p.root,
		substitutions: p.substitutions,
	}
}

func (p *Config) GetNode(path string) (*hocon.HoconValue, error) {
	if p == nil {
		return nil, nil
	}

	elements := splitDottedPathHonouringQuotes(path)
	currentNode := p.root

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
			if p.fallback != nil {
				return p.fallback.GetNode(path)
			}
			return nil, nil
		}
	}
	return currentNode, nil
}

func (p *Config) GetBoolean(path string, defaultVal ...bool) (bool, error) {
	obj, err := p.GetNode(path)
	if err != nil {
		return false, err
	}

	if obj == nil {
		if len(defaultVal) > 0 {
			return defaultVal[0], nil
		}
		return false, nil
	}
	return obj.GetBoolean()
}

func (p *Config) GetByteSize(path string) (*big.Int, error) {
	obj, err := p.GetNode(path)
	if err != nil {
		return nil, err
	}

	if obj == nil {
		return big.NewInt(-1), nil
	}
	return obj.GetByteSize()
}

func (p *Config) GetInt32(path string, defaultVal ...int32) (int32, error) {
	obj, err := p.GetNode(path)
	if err != nil {
		return 0, err
	}

	if obj == nil {
		if len(defaultVal) > 0 {
			return defaultVal[0], nil
		}
		return 0, nil
	}
	return obj.GetInt32()
}

func (p *Config) GetInt64(path string, defaultVal ...int64) (int64, error) {
	obj, err := p.GetNode(path)
	if err != nil {
		return 0, err
	}

	if obj == nil {
		if len(defaultVal) > 0 {
			return defaultVal[0], nil
		}
		return 0, nil
	}
	return obj.GetInt64()
}

func (p *Config) GetString(path string, defaultVal ...string) (string, error) {
	obj, err := p.GetNode(path)
	if err != nil {
		return "", err
	}

	if obj == nil {
		if len(defaultVal) > 0 {
			return defaultVal[0], nil
		}
		return "", nil
	}
	return obj.GetString()
}

func (p *Config) GetFloat32(path string, defaultVal ...float32) (float32, error) {
	obj, err := p.GetNode(path)
	if err != nil {
		return 0, err
	}

	if obj == nil {
		if len(defaultVal) > 0 {
			return defaultVal[0], nil
		}
		return 0, nil
	}
	return obj.GetFloat32()
}

func (p *Config) GetFloat64(path string, defaultVal ...float64) (float64, error) {
	obj, err := p.GetNode(path)
	if err != nil {
		return 0, err
	}

	if obj == nil {
		if len(defaultVal) > 0 {
			return defaultVal[0], nil
		}
		return 0, nil
	}
	return obj.GetFloat64()
}

func (p *Config) GetTimeDuration(path string, defaultVal ...time.Duration) (time.Duration, error) {
	obj, err := p.GetNode(path)
	if err != nil {
		return 0, err
	}

	if obj == nil {
		if len(defaultVal) > 0 {
			return defaultVal[0], nil
		}
		return 0, nil
	}
	return obj.GetTimeDuration(true)
}

func (p *Config) GetTimeDurationInfiniteNotAllowed(path string, defaultVal ...time.Duration) (time.Duration, error) {
	obj, err := p.GetNode(path)
	if err != nil {
		return 0, err
	}

	if obj == nil {
		if len(defaultVal) > 0 {
			return defaultVal[0], nil
		}
		return 0, nil
	}
	return obj.GetTimeDuration(false)
}

func (p *Config) GetBooleanList(path string) ([]bool, error) {
	obj, err := p.GetNode(path)
	if err != nil {
		return nil, err
	}

	if obj == nil {
		return nil, nil
	}
	return obj.GetBooleanList()
}

func (p *Config) GetFloat32List(path string) ([]float32, error) {
	obj, err := p.GetNode(path)
	if err != nil {
		return nil, err
	}

	if obj == nil {
		return nil, nil
	}
	return obj.GetFloat32List()
}

func (p *Config) GetFloat64List(path string) ([]float64, error) {
	obj, err := p.GetNode(path)
	if err != nil {
		return nil, err
	}

	if obj == nil {
		return nil, nil
	}
	return obj.GetFloat64List()
}

func (p *Config) GetInt32List(path string) ([]int32, error) {
	obj, err := p.GetNode(path)
	if err != nil {
		return nil, err
	}

	if obj == nil {
		return nil, nil
	}
	return obj.GetInt32List()
}

func (p *Config) GetInt64List(path string) ([]int64, error) {
	obj, err := p.GetNode(path)
	if err != nil {
		return nil, err
	}

	if obj == nil {
		return nil, nil
	}
	return obj.GetInt64List()
}

func (p *Config) GetByteList(path string) ([]byte, error) {
	obj, err := p.GetNode(path)
	if err != nil {
		return nil, err
	}

	if obj == nil {
		return nil, nil
	}
	return obj.GetByteList()
}

func (p *Config) GetStringList(path string) ([]string, error) {
	obj, err := p.GetNode(path)
	if err != nil {
		return nil, err
	}

	if obj == nil {
		return nil, nil
	}
	return obj.GetStringList()
}

func (p *Config) GetConfig(path string) (*Config, error) {
	if p == nil {
		return nil, nil
	}

	value, err := p.GetNode(path)
	if err != nil {
		return nil, err
	}

	if p.fallback != nil {
		f, err := p.fallback.GetConfig(path)
		if err != nil {
			return nil, err
		}

		if value == nil && f == nil {
			return nil, nil
		}
		if value == nil {
			return f, nil
		}
		root, err := NewConfigFromRoot(hocon.NewHoconRoot(value))
		if err != nil {
			return nil, err
		}

		return root.WithFallback(f)
	}

	if value == nil {
		return nil, nil
	}
	return NewConfigFromRoot(hocon.NewHoconRoot(value))
}

func (p *Config) GetValue(path string) (*hocon.HoconValue, error) {
	return p.GetNode(path)
}

func (p *Config) WithFallback(fallback *Config) (*Config, error) {
	if fallback == p {
		return nil, nil
	}

	if fallback == nil {
		return p, nil
	}

	selfObjectV, err := p.root.GetObject()
	if err != nil {
		return nil, err
	}

	fallbackObjectV, err := fallback.root.GetObject()
	if err != nil {
		return nil, err
	}

	mergedRoot, err := selfObjectV.MergeImmutable(fallbackObjectV)
	if err != nil {
		return nil, err
	}

	newRoot := hocon.NewHoconValue()

	newRoot.AppendValue(mergedRoot)

	mergedConfig := p.Copy(fallback)

	mergedConfig.root = newRoot

	return mergedConfig, nil
}

func (p *Config) HasPath(path string) bool {
	node, err := p.GetNode(path)
	if err != nil {
		return false
	}

	return node != nil
}

func (p *Config) IsObject(path string) (bool, error) {
	node, err := p.GetNode(path)
	if err != nil || node == nil {
		return false, nil
	}

	return node.IsObject()
}

func (p *Config) IsArray(path string) bool {
	node, err := p.GetNode(path)
	if err != nil || node == nil {
		return false
	}

	return node.IsArray()
}

func (p *Config) AddConfig(textConfig string, fallbackConfig *Config) (*Config, error) {
	root, err := hocon.Parse(textConfig, nil)
	if err != nil {
		return nil, err
	}

	config, err := NewConfigFromRoot(root)
	if err != nil {
		return nil, err
	}

	return config.WithFallback(fallbackConfig)
}

func (p *Config) AddConfigWithTextFallback(config *Config, textFallback string) (*Config, error) {
	fallbackRoot, err := hocon.Parse(textFallback, nil)
	if err != nil {
		return nil, err
	}

	fallbackConfig, err := NewConfigFromRoot(fallbackRoot)
	if err != nil {
		return nil, err
	}

	return config.WithFallback(fallbackConfig)
}

func (p Config) String() string {
	return p.root.String()
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
