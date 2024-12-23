package telegrammsgsender

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	jsoniter "github.com/json-iterator/go"
)

var (
	senders sync.Map
	json    = jsoniter.ConfigCompatibleWithStandardLibrary
	//errChatId   string
	debugChatId   string
	debugBotToken string
	debugThreadId int
)

type sender struct {
	BotToken string
	Queue    []msgPayload
	Mu       *sync.Mutex
	Waiter   *sync.WaitGroup
	Waiting  bool
}

type msgPayload struct {
	Type        string
	Bytes       []byte
	ContentType string
	IsError     bool
}

func (msg Message) Send(botToken string, mergeEmbeds ...bool) error {
	sender := getSender(botToken)
	return sender.queueAdd(msg, false)
}

func (msg Message) SendV2(botToken string, chatId string, threadId int, mergeEmbeds ...bool) error {

	if debugBotToken != "" {
		botToken = debugBotToken
		chatId = debugChatId
		threadId = debugThreadId
	}
	msg.ChatID = chatId
	msg.MessageThreadID = threadId
	sender := getSender(botToken)
	return sender.queueAdd(msg, false)
}

// Set global error webhook url
// for unset, just set to empty string
// func SetErrorChatId(chatId string) {
// 	errChatId = chatId
// }

// Set debug webhook
// that override every whs
func SetDebugBot(botToken string, chatId string, threadId int) {
	debugBotToken = botToken
	debugChatId = chatId
	debugThreadId = threadId
}

func newSender(token string) *sender {
	return &sender{
		BotToken: token,
		Queue:    []msgPayload{},
		Mu:       &sync.Mutex{},
		Waiter:   &sync.WaitGroup{},
		Waiting:  false,
	}
}

func getSender(token string) *sender {
	s, found := senders.LoadOrStore(token, newSender(token))
	sender := s.(*sender)
	if !found {
		sender.initSender()
	}
	return sender
}

func (s *sender) initSender() {
	go func() {
		for {
			s.Waiter.Wait()
			if p := s.queueGet(); len(p.Bytes) > 0 {
				retry := true
				for retry {
					res, err := http.Post(
						fmt.Sprintf("https://api.telegram.org/bot%s/%s", s.BotToken, p.Type),
						p.ContentType,
						bytes.NewBuffer(p.Bytes),
					)
					if err != nil {
						continue
					}

					switch res.StatusCode {
					case 200, 204:
						//rtRemaining, _ := strconv.Atoi(res.Header.Get("x-ratelimit-remaining"))
						//if rtRemaining < 3 {
						time.Sleep(300 * time.Millisecond)
						//}
						retry = false
					case 429:
						//ratelimitDelay, _ := strconv.Atoi(res.Header.Get("retry-after"))
						log.Printf("[telegram-messagesender] Ratelimited - %s\n", s.BotToken)
						time.Sleep(time.Duration(3000) * time.Millisecond)
						retry = true
					default:
						bbody, _ := io.ReadAll(res.Body)
						log.Printf("[telegram-messagesender] SendError(%d): - %s\n", res.StatusCode, string(bbody))
						log.Printf("[telegram-messagesender] Payload: %s", string(p.Bytes))
						retry = false
					}
					res.Body.Close()
				}
			}
		}
	}()
}
