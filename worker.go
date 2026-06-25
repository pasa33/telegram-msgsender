package telegrammsgsender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var (
	workers   sync.Map
	workersMu sync.Mutex
)

type worker struct {
	token string
	queue chan payload
}

type payload struct {
	method      string
	body        []byte
	contentType string
}

func getWorker(token string) *worker {
	if v, ok := workers.Load(token); ok {
		return v.(*worker)
	}
	workersMu.Lock()
	defer workersMu.Unlock()
	if v, ok := workers.Load(token); ok {
		return v.(*worker)
	}
	w := &worker{
		token: token,
		queue: make(chan payload, 100),
	}
	workers.Store(token, w)
	go w.run()
	return w
}

func (w *worker) enqueue(msg *Message, chatID string, threadID int) error {
	p, err := buildPayload(msg, chatID, threadID)
	if err != nil {
		return err
	}
	select {
	case w.queue <- p:
		return nil
	default:
		return fmt.Errorf("telegram-msgsender: queue full (100 messages pending)")
	}
}

func (w *worker) run() {
	for p := range w.queue {
		w.send(p)
	}
}

func (w *worker) send(p payload) {
	for {
		resp, err := http.Post(
			fmt.Sprintf("https://api.telegram.org/bot%s/%s", w.token, p.method),
			p.contentType,
			bytes.NewReader(p.body),
		)
		if err != nil {
			log.Printf("[telegram-msgsender] HTTP error: %v — retrying in 5s", err)
			time.Sleep(5 * time.Second)
			continue
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		switch resp.StatusCode {
		case 200, 201:
			time.Sleep(50 * time.Millisecond)
			return
		case 429:
			delay := retryAfter(body, resp.Header)
			log.Printf("[telegram-msgsender] rate limited, retrying in %v", delay)
			time.Sleep(delay)
		default:
			log.Printf("[telegram-msgsender] error %d: %s", resp.StatusCode, body)
			return
		}
	}
}

func retryAfter(body []byte, h http.Header) time.Duration {
	var r struct {
		Parameters struct {
			RetryAfter int `json:"retry_after"`
		} `json:"parameters"`
	}
	if err := json.Unmarshal(body, &r); err == nil && r.Parameters.RetryAfter > 0 {
		return time.Duration(r.Parameters.RetryAfter) * time.Second
	}
	if s := h.Get("Retry-After"); s != "" {
		if n, err := strconv.Atoi(s); err == nil {
			return time.Duration(n) * time.Second
		}
	}
	return 3 * time.Second
}

func buildPayload(msg *Message, chatID string, threadID int) (payload, error) {
	if len(msg.imageBytes) > 0 {
		return buildPhotoMultipart(msg, chatID, threadID)
	}
	if msg.imageURL != "" {
		return buildPhotoJSON(msg, chatID, threadID)
	}
	return buildMessageJSON(msg, chatID, threadID)
}

func buildMessageJSON(msg *Message, chatID string, threadID int) (payload, error) {
	m := &sendMessage{
		ChatID:             chatID,
		Text:               msg.text(),
		ParseMode:          msg.parseMode,
		ProtectContent:     msg.protect,
		DisableLinkPreview: msg.noPreview,
	}
	if threadID != 0 {
		m.MessageThreadID = &threadID
	}
	data, err := json.Marshal(m)
	if err != nil {
		return payload{}, err
	}
	return payload{method: "sendMessage", body: data, contentType: "application/json"}, nil
}

func buildPhotoJSON(msg *Message, chatID string, threadID int) (payload, error) {
	p := &sendPhoto{
		ChatID:         chatID,
		Photo:          msg.imageURL,
		Caption:        msg.text(),
		ParseMode:      msg.parseMode,
		ProtectContent: msg.protect,
	}
	if threadID != 0 {
		p.MessageThreadID = &threadID
	}
	data, err := json.Marshal(p)
	if err != nil {
		return payload{}, err
	}
	return payload{method: "sendPhoto", body: data, contentType: "application/json"}, nil
}

func buildPhotoMultipart(msg *Message, chatID string, threadID int) (payload, error) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	fw, err := w.CreateFormFile("photo", "photo")
	if err != nil {
		return payload{}, err
	}
	if _, err = fw.Write(msg.imageBytes); err != nil {
		return payload{}, err
	}

	fields := [][2]string{
		{"chat_id", chatID},
		{"caption", msg.text()},
		{"parse_mode", msg.parseMode},
		{"protect_content", strconv.FormatBool(msg.protect)},
	}
	if threadID != 0 {
		fields = append(fields, [2]string{"message_thread_id", strconv.Itoa(threadID)})
	}
	for _, f := range fields {
		if err = w.WriteField(f[0], f[1]); err != nil {
			return payload{}, err
		}
	}
	w.Close()

	return payload{method: "sendPhoto", body: buf.Bytes(), contentType: w.FormDataContentType()}, nil
}
