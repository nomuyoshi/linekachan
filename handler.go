package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"
)

// PREFIX はリマインド登録に必要なキーワード
const PREFIX = "リマインド"

// USAGE は使い方説明。
const USAGE = `「リマインド　内容」と話しかけてください。
例) リマインド　食パンを買って帰る`

// CallbackHandler はLINEからのcallbackを処理する構造体
type CallbackHandler struct {
	bot *linebot.Client
}

func (h *CallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	events, err := h.bot.ParseRequest(r)
	if err != nil {
		log.Print("Error parse request:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			var messages []linebot.SendingMessage

			switch receivedMessage := event.Message.(type) {
			case *linebot.TextMessage:
				receivedText := strings.TrimSpace(receivedMessage.Text)

				if strings.HasPrefix(receivedText, PREFIX) {
					content := strings.TrimSpace(strings.Replace(receivedText, PREFIX, "", 1))
					messages = append(messages, linebot.NewTextMessage("受け付けました！！\n"+content))
				} else {
					messages = append(messages, linebot.NewTextMessage(USAGE))
				}
			default:
				messages = append(messages, linebot.NewTextMessage(USAGE))
			}

			if _, err := h.bot.ReplyMessage(event.ReplyToken, messages...).Do(); err != nil {
				log.Print("Error reply: ", err)
			}
		}
	}
	w.WriteHeader(http.StatusOK)
}
