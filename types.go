package telegrammsgsender

type sendMessage struct {
	ChatID             string `json:"chat_id"`
	MessageThreadID    *int   `json:"message_thread_id,omitempty"`
	Text               string `json:"text"`
	ParseMode          string `json:"parse_mode"`
	ProtectContent     bool   `json:"protect_content"`
	DisableLinkPreview bool   `json:"disable_web_page_preview"`
}

type sendPhoto struct {
	ChatID          string `json:"chat_id"`
	MessageThreadID *int   `json:"message_thread_id,omitempty"`
	Photo           string `json:"photo"`
	Caption         string `json:"caption"`
	ParseMode       string `json:"parse_mode"`
	ProtectContent  bool   `json:"protect_content"`
}
