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
	return ValidateLine(text)
}

func MakeUrl(text, url string) string {
	return fmt.Sprintf("[%s](%s)", ValidateLine(text), ValidateUrl(url))
}

func ValidateUrl(s string) string {
	for _, v := range urlChar {
		s = strings.ReplaceAll(s, v, `\`+v)
	}
	return s
}

func ValidateLine(s string) string {
	for _, v := range txtChar {
		s = strings.ReplaceAll(s, v, `\`+v)
	}
	return s
}
