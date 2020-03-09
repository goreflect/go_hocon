package hocon

import "strings"

const newLine = "\r\n"

// quoteStringIfNeeded wrap text with quotes if it contains a space or tab symbol
func quoteStringIfNeeded(text string) string {
	if strings.IndexByte(text, ' ') >= 0 ||
		strings.IndexByte(text, '\t') >= 0 {
		return "\"" + text + "\""
	}
	return text
}
