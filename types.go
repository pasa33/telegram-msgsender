package telegrammsgsender

//cusotm wrapper
type Message struct {
	ChatID          string
	MessageThreadID int
	ContentLines    []string
	Image           string
	ImageBytes      []byte
	ParseMode       string
	ProtectContent  bool
}

type sendPhoto struct {
	ChatID          string `json:"chat_id"`
	MessageThreadID *int   `json:"message_thread_id,omitempty"`
	Photo           string `json:"photo"`
	Caption         string `json:"caption"`
	ParseMode       string `json:"parse_mode"`
	ProtectContent  bool   `json:"protect_content"`
	ImgBytes        []byte `json:"-"`
}

type sendMessage struct {
	ChatID          string `json:"chat_id"`
	MessageThreadID *int   `json:"message_thread_id,omitempty"`
	Text            string `json:"text"`
	ParseMode       string `json:"parse_mode"`
	ProtectContent  bool   `json:"protect_content"`
}
