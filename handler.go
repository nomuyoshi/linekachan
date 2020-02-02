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
const lineLayout = "2006-01-02T15:04"
const lineMessageLayout = "2006/01/02 15:04"

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
		var messages []linebot.SendingMessage

		switch event.Type {
		case linebot.EventTypeMessage:
			if receivedMessage, ok := event.Message.(*linebot.TextMessage); ok {
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
			} else {
				messages = append(messages, linebot.NewTextMessage(usage))
			}
		case linebot.EventTypePostback:
			if !strings.HasPrefix(event.Postback.Data, "scheduleId=") {
				continue
			}
			scheduleID := strings.TrimPrefix(event.Postback.Data, "scheduleId=")

			if schedule, _ := findScheduleBy(scheduleID, event.Source.UserID); schedule != nil {
				datetime, _ := time.Parse(lineLayout, event.Postback.Params.Datetime)
				schedule.update(datetime)
				messages = append(messages, linebot.NewTextMessage(datetime.Format(lineMessageLayout)+"で受け付けました！"))
			} else {
				messages = append(messages, linebot.NewTextMessage("予期せぬエラーが発生しました。"))
			}
		}

		if _, err := h.bot.ReplyMessage(event.ReplyToken, messages...).Do(); err != nil {
			log.Print("Error reply: ", err)
		}
	}
	w.WriteHeader(http.StatusOK)
}

func buildDatetimePickerMessage(postback string) *linebot.TemplateMessage {
	now := time.Now()
	max := now.AddDate(1, 0, 0)
	datetimePickerAction := linebot.NewDatetimePickerAction(
		datetimePickerLabel,
		postback,
		"datetime",
		now.Format(lineLayout),
		max.Format(lineLayout),
		now.Format(lineLayout),
	)
	datetimeTemplate := linebot.NewButtonsTemplate("", "", datetimeSelectText, datetimePickerAction)

	return linebot.NewTemplateMessage("Lineをアップデートしてください", datetimeTemplate)
}
