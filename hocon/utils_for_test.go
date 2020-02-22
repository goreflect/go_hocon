package hocon

const (
	simpleKey1 = "key1"
	simpleKey2 = "key2"
	simpleKey3 = "key3"
	specials   = "`-=~!@#$%^&*()_+[]\\{}|\"':;,./<>?"
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

	simpleLiteral1 = NewHoconLiteral("value1")
	simpleLiteral2 = NewHoconLiteral("value2")
	simpleLiteral3 = NewHoconLiteral("value3")
)

func getMapOfTwoSimpleLiterals() map[string]*HoconValue {
	return map[string]*HoconValue{
		simpleKey1: wrapInValue(simpleLiteral1),
		simpleKey2: wrapInValue(simpleLiteral2),
	}
}

func getArrayOfTwoSimpleKeys() []string {
	return []string{simpleKey1, simpleKey2}
}

func getCycledSubstitution() *HoconSubstitution {
	cycledSubstitution := &HoconSubstitution{}
	cycledSubstitution.ResolvedValue = &HoconValue{values: []HoconElement{cycledSubstitution}}
	return cycledSubstitution
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

// wrapAllInObject wraps gotten values into HoconObject. It was made for test purposes
// only. Make sure that size of elements are equal!
func wrapAllInObject(keys []string, elements []HoconElement) *HoconObject {
	items := map[string]*HoconValue{}
	for i, key := range keys {
		//noinspection GoNilness
		items[key] = wrapInValue(elements[i])
	}
	return &HoconObject{
		keys:  keys,
		items: items,
	}
}

func wrapInValue(objects ...HoconElement) *HoconValue {
	return &HoconValue{values: objects}
}

func wrapInSubstitution(value HoconElement) *HoconSubstitution {
	return &HoconSubstitution{
		Path:          "",
		ResolvedValue: wrapInValue(value),
		IsOptional:    false,
		OriginalPath:  "",
	}

}

func wrapInArray(values ...*HoconValue) *HoconArray {
	return &HoconArray{values: values}
}
