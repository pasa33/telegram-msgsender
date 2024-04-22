package telegrammsgsender

import (
	"bytes"
	"cmp"
	"fmt"
	"mime/multipart"
	"strings"
)

func (s *sender) queueGet() (p msgPayload) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	if len(s.Queue) > 0 {
		p = s.Queue[0]
		if len(s.Queue) > 1 {
			s.Queue = s.Queue[1:]
		} else {
			s.Queue = []msgPayload{}
		}
		return
	}
	s.Waiting = true
	s.Waiter.Add(1)
	return
}

func (s *sender) queueAdd(msg Message, isErr bool) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	p := msgPayload{
		IsError: isErr,
	}

	if msg.Image != "" || len(msg.ImageBytes) > 0 {

		sP := sendPhoto{
			ChatID:         cmp.Or(debugChatId, msg.ChatID),
			Photo:          msg.Image,
			Caption:        strings.Join(msg.ContentLines, "\n"),
			ParseMode:      cmp.Or(msg.ParseMode, "MarkdownV2"),
			ProtectContent: msg.ProtectContent,
		}
		if msg.MessageThreadID != 0 && debugChatId == "" {
			sP.MessageThreadID = &msg.MessageThreadID
		}

		if len(msg.ImageBytes) > 0 {
			buffer := new(bytes.Buffer)
			writer := multipart.NewWriter(buffer)

			fw, err := writer.CreateFormFile("photo", "photo")
			if err != nil {
				return err
			}

			if _, err := fw.Write(msg.ImageBytes); err != nil {
				return err
			}

			if err := writer.WriteField("chat_id", sP.ChatID); err != nil {
				return err
			}
			if sP.MessageThreadID != nil {
				if err := writer.WriteField("message_thread_id", fmt.Sprintln(*sP.MessageThreadID)); err != nil {
					return err
				}
			}
			if err := writer.WriteField("photo", "attach://photo"); err != nil {
				return err
			}
			if err := writer.WriteField("caption", sP.Caption); err != nil {
				return err
			}
			if err := writer.WriteField("parse_mode", sP.ParseMode); err != nil {
				return err
			}
			if err := writer.WriteField("protect_content", "true"); err != nil {
				return err
			}
			writer.Close()

			p.Bytes = buffer.Bytes()
			p.Type = "sendPhoto"
			p.ContentType = writer.FormDataContentType()

		} else {
			data, err := json.Marshal(sP)
			if err != nil {
				return err
			}

			p.Bytes = data
			p.Type = "sendPhoto"
			p.ContentType = "application/json"
		}
	} else { //without photo
		sM := sendMessage{
			ChatID:         cmp.Or(debugChatId, msg.ChatID),
			Text:           strings.Join(msg.ContentLines, "\n"),
			ParseMode:      cmp.Or(msg.ParseMode, "MarkdownV2"),
			ProtectContent: msg.ProtectContent,
		}
		if msg.MessageThreadID != 0 && debugChatId == "" {
			sM.MessageThreadID = &msg.MessageThreadID
		}

		data, err := json.Marshal(sM)
		if err != nil {
			return err
		}

		p.Bytes = data
		p.Type = "sendMessage"
		p.ContentType = "application/json"
	}

	if isErr {
		s.Queue = append([]msgPayload{p}, s.Queue...)
	} else {
		s.Queue = append(s.Queue, p)
	}
	if s.Waiting {
		s.Waiting = false
		s.Waiter.Done()
	}
	return nil
}
