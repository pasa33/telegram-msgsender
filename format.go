package telegrammsgsender

import (
	"fmt"
	"strings"
)

var mdV2Chars = []string{`\`, "_", "*", "[", "]", "(", ")", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!"}
var mdV2URLChars = []string{`\`, ")"}

// Escape escapes all MarkdownV2 special characters in a plain text string.
func Escape(s string) string {
	for _, c := range mdV2Chars {
		s = strings.ReplaceAll(s, c, `\`+c)
	}
	return s
}

// Bold returns the text as MarkdownV2 bold.
func Bold(text string) string {
	return fmt.Sprintf("*%s*", Escape(text))
}

// Italic returns the text as MarkdownV2 italic.
func Italic(text string) string {
	return fmt.Sprintf("_%s_", Escape(text))
}

// Code returns the text as MarkdownV2 inline code.
// Only \ and ` are escaped inside code spans, per Telegram spec.
func Code(text string) string {
	text = strings.ReplaceAll(text, `\`, `\\`)
	text = strings.ReplaceAll(text, "`", "\\`")
	return fmt.Sprintf("`%s`", text)
}

// Pre returns the text as a MarkdownV2 pre-formatted block.
// Only \ and ` are escaped inside pre blocks, per Telegram spec.
func Pre(text string) string {
	text = strings.ReplaceAll(text, `\`, `\\`)
	text = strings.ReplaceAll(text, "`", "\\`")
	return fmt.Sprintf("```\n%s\n```", text)
}

// Link returns a MarkdownV2 hyperlink.
func Link(label, url string) string {
	for _, c := range mdV2URLChars {
		url = strings.ReplaceAll(url, c, `\`+c)
	}
	return fmt.Sprintf("[%s](%s)", Escape(label), url)
}

// Field returns a formatted "key: value" line, like a Discord embed field.
func Field(name, value string) string {
	return fmt.Sprintf("*%s*: %s", Escape(name), Escape(value))
}
