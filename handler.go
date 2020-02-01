package main

import (
	"log"
	"net/http"

	"github.com/line/line-bot-sdk-go/linebot"
)

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
			replyToken := event.ReplyToken
			var messages []linebot.SendingMessage
			message := linebot.NewTextMessage("受け付けました！！")
			messages = append(messages, message)

			if _, err := h.bot.ReplyMessage(replyToken, messages...).Do(); err != nil {
				log.Print("Error reply: ", err)
			}
		}
	}
	w.WriteHeader(http.StatusOK)
}
