package telegrammsgsender

import "strings"

// Message is assembled using chained methods, similar to a Discord embed.
type Message struct {
	lines      []string
	imageURL   string
	imageBytes []byte
	parseMode  string
	protect    bool
	noPreview  bool
}

// NewMessage creates an empty Message with MarkdownV2 as the default parse mode.
func NewMessage() *Message {
	return &Message{parseMode: "MarkdownV2"}
}

// Line appends a line. Use Bold(), Italic(), Field() etc. to format inline.
func (m *Message) Line(text string) *Message {
	m.lines = append(m.lines, text)
	return m
}

// Blank appends an empty line.
func (m *Message) Blank() *Message {
	m.lines = append(m.lines, "")
	return m
}

// Image attaches an image by URL.
func (m *Message) Image(url string) *Message {
	m.imageURL = url
	return m
}

// ImageBytes attaches an image from raw bytes. Prefer this over Image() when the
// URL may be banned by Telegram (common with scraped images).
func (m *Message) ImageBytes(b []byte) *Message {
	m.imageBytes = b
	return m
}

// HTML switches the parse mode to HTML. Bold/Italic/Code helpers still work;
// Escape() is not needed since HTML entities are different.
func (m *Message) HTML() *Message {
	m.parseMode = "HTML"
	return m
}

// Protect disables forwarding and saving of the message.
func (m *Message) Protect() *Message {
	m.protect = true
	return m
}

// NoPreview disables link preview for text messages.
func (m *Message) NoPreview() *Message {
	m.noPreview = true
	return m
}

func (m *Message) text() string {
	return strings.Join(m.lines, "\n")
}
