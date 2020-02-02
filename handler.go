package main

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
)

// リマインド登録に必要なキーワード
const prefix = "リマインド"
const usage = `「リマインド　内容」と話しかけてください。
例) リマインド　食パンを買って帰る`
const datetimeSelectText = "リマインドして欲しい日時を指定してください。"
const datetimePickerLabel = "日時指定"

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

				if strings.HasPrefix(receivedText, prefix) {
					content := strings.Replace(receivedText, prefix, "", 1)
					schedule := newSchedule(event.Source.UserID, content)
					if err := schedule.create(); err != nil {
						break
					}

					messages = append(messages, buildDatetimePickerMessage(schedule.postbackData()))
				} else {
					messages = append(messages, linebot.NewTextMessage(usage))
				}
			default:
				messages = append(messages, linebot.NewTextMessage(usage))
			}

			if _, err := h.bot.ReplyMessage(event.ReplyToken, messages...).Do(); err != nil {
				log.Print("Error reply: ", err)
			}
		}
	}
	w.WriteHeader(http.StatusOK)
}

func buildDatetimePickerMessage(postback string) *linebot.TemplateMessage {
	now := time.Now()
	max := now.AddDate(1, 0, 0)
	layout := "2006-01-02T15:04"
	datetimePickerAction := linebot.NewDatetimePickerAction(
		datetimePickerLabel,
		postback,
		"datetime",
		now.Format(layout),
		max.Format(layout),
		now.Format(layout),
	)
	datetimeTemplate := linebot.NewButtonsTemplate("", "", datetimeSelectText, datetimePickerAction)

	return linebot.NewTemplateMessage("Lineをアップデートしてください", datetimeTemplate)
}
