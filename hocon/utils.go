package hocon

import "strings"

// quoteStringIfNeeded wrap text with quotes if it contains a space or tab symbol
func quoteStringIfNeeded(text string) string {
	if strings.IndexByte(text, ' ') >= 0 ||
		strings.IndexByte(text, '\t') >= 0 {
		return "\"" + text + "\""
	}
	return text
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
