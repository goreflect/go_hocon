package configuration

import (
	"encoding/json"
	"io/ioutil"

	"github.com/goreflect/go_hocon/hocon"
)

func ParseString(text string, includeCallback ...hocon.IncludeCallback) (*Config, error) {
	var callback hocon.IncludeCallback
	if len(includeCallback) > 0 {
		callback = includeCallback[0]
	} else {
		callback = defaultIncludeCallback
	}
	root, err := hocon.Parse(text, callback)
	if err != nil {
		return nil, err
	}

	return NewConfigFromRoot(root)
}

func LoadConfig(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return ParseString(string(data), defaultIncludeCallback)
}

func FromObject(obj interface{}) (*Config, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	return ParseString(string(data), defaultIncludeCallback)
}

func defaultIncludeCallback(filename string) (*hocon.HoconRoot, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return hocon.Parse(string(data), defaultIncludeCallback)
}
