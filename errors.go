package telegrammsgsender

// func sendError(token string, status string, req, res []byte) error {
// 	sender := getSender(token)
// 	//sender := getSender(cmp.Or(debugUrl, errUrl, url))

// 	return sender.queueAdd(makeErrorMsg(status, url, req, res), true)
// }

// func makeErrorMsg(status, url string, req, res []byte) Message {
// 	c := status
// 	if len(errUrl) > 0 {
// 		c += fmt.Sprintf("\n`%s`", url)
// 	}
// 	return Message{
// 		ChatID:  cmp.Or(debugChatId, errChatId, ),
// 		Content: c,
// 		Files: []File{
// 			{Name: "ReqPayload.txt", Bytes: []byte(base64.StdEncoding.EncodeToString(req))},
// 			{Name: "ResPayload.txt", Bytes: []byte(base64.StdEncoding.EncodeToString(res))},
// 		},
// 	}
// }
