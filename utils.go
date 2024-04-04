package telegrammsgsender

import (
	"fmt"
	"strings"
)

var txtChar = []string{`\`, "_", "*", "[", "]", "(", ")", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!"}
var urlChar = []string{`\`, ")"}

func (m *Message) AddLine(s string) {
	m.ContentLines = append(m.ContentLines, s)
}

func MakeText(text string) string {
	return validateLine(text)
}

func MakeUrl(text, url string) string {
	return fmt.Sprintf("[%s](%s)", validateLine(text), validateUrl(url))
}

func validateUrl(s string) string {
	for _, v := range urlChar {
		s = strings.ReplaceAll(s, v, `\`+v)
	}
	return s
}

func validateLine(s string) string {
	for _, v := range txtChar {
		s = strings.ReplaceAll(s, v, `\`+v)
	}
	return s
}
