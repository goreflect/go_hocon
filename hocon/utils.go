package hocon

import (
	"strings"
)

const (
	newLine string = "\r\n"
)

var (
	simpleTwoValuesArray = wrapInArray(wrapInValue(NewHoconLiteral("a")), wrapInValue(NewHoconLiteral("b")))

	simpleArrayWithOldValue = &HoconArray{values: []*HoconValue{
		{
			values:   []HoconElement{NewHoconLiteral("current")},
			oldValue: wrapInValue(NewHoconLiteral("old")),
		},
	}}

	simpleObject = makeHoconObject([]string{"key"}, []string{"value"})
	//simpleObject = &HoconObject{
	//	keys:  []string{"a"},
	//	items: map[string]*HoconValue{"a": {values: []HoconElement{NewHoconLiteral("b")}}},
	//}

	simpleLiteral = NewHoconLiteral("a")

	simpleNestedObject = &HoconObject{
		keys:  []string{"a"},
		items: map[string]*HoconValue{"a": {values: []HoconElement{simpleObject}}},
	}
)

// quoteStringIfNeeded wrap text with quotes if it contains a space or tab symbol
func quoteStringIfNeeded(text string) string {
	if strings.IndexByte(text, ' ') >= 0 ||
		strings.IndexByte(text, '\t') >= 0 {
		return "\"" + text + "\""
	}
	return text
}

func getCycledSubstitution() *HoconSubstitution {
	cycledSubstitution := &HoconSubstitution{}
	cycledSubstitution.ResolvedValue = &HoconValue{values: []HoconElement{cycledSubstitution}}
	return cycledSubstitution
}

func getCycledObject() *HoconObject {
	return &HoconObject{
		keys:  []string{"a"},
		items: map[string]*HoconValue{"a": {values: []HoconElement{getCycledSubstitution()}}},
	}
}

func getCycledSubstitutionValue() *HoconValue {
	return &HoconValue{values: []HoconElement{getCycledSubstitution()}}
}

func getCycledObjectValue() *HoconValue {
	return &HoconValue{values: []HoconElement{getCycledObject()}}
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

func wrapInObject(key string, element HoconElement) *HoconObject {
	return &HoconObject{
		keys:  []string{key},
		items: map[string]*HoconValue{key: wrapInValue(element)},
	}
}

func wrapInValue(object HoconElement) *HoconValue {
	return &HoconValue{values: []HoconElement{object}}
}

func wrapInSubstitution(value *HoconValue) *HoconSubstitution {
	return &HoconSubstitution{
		Path:          "",
		ResolvedValue: value,
		IsOptional:    false,
		OriginalPath:  "",
	}

}

func wrapInArray(values ...*HoconValue) *HoconArray {
	return &HoconArray{values: values}
}
