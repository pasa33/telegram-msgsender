package telegrammsgsender

var debugCfg *debugConfig

type debugConfig struct {
	BotToken string
	ChatID   string
	ThreadID int
}

// Send queues a message for async delivery. threadID can be 0 if not needed.
func Send(botToken, chatID string, threadID int, msg *Message) error {
	token, cid, tid := botToken, chatID, threadID
	if debugCfg != nil {
		token = debugCfg.BotToken
		cid = debugCfg.ChatID
		tid = debugCfg.ThreadID
	}
	return getWorker(token).enqueue(msg, cid, tid)
}

// SetDebug overrides all sends to go to the specified bot/chat. Call with empty botToken to clear.
func SetDebug(botToken, chatID string, threadID int) {
	if botToken == "" {
		debugCfg = nil
		return
	}
	debugCfg = &debugConfig{BotToken: botToken, ChatID: chatID, ThreadID: threadID}
}
