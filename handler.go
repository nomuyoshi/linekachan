package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
)

// リマインド登録に必要なキーワード
const (
	prefix      = "リマインド"
	usage       = "「リマインド　内容」と話しかけて\u270B\n例) リマインド　食パンを買って帰る"
	pickerText  = "いつ言えばいい\u2753"
	pickerLabel = "日時指定\U0001F5D3"
	okText      = "了解\U0001F44C\u203C"
	errorText   = "エラーが発生したみたい...\U0001F605\nもう一度やり直してみて\U0001F64F"
	lineLayout  = "2006-01-02T15:04"
	strLayout   = "2006/01/02 15:04"
)

// CallbackHandler はLINEからのcallbackを処理する構造体
type CallbackHandler struct {
	bot  *linebot.Client
	lkDb *LineKachanDb
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
			if resMessage, ok := event.Message.(*linebot.TextMessage); ok {
				resText := strings.TrimSpace(resMessage.Text)

				if strings.HasPrefix(resText, prefix) {
					schedule := NewSchedule(event.Source.UserID, strings.Replace(resText, prefix, "", 1))
					if err := h.lkDb.AddSchedule(schedule); err != nil {
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
			scheduleID, _ := strconv.ParseInt(strings.TrimPrefix(event.Postback.Data, "scheduleId="), 10, 64)
			schedule, err := h.lkDb.FindScheduleBy(scheduleID, event.Source.UserID)
			if err == nil {
				datetime, _ := time.ParseInLocation(lineLayout, event.Postback.Params.Datetime, time.Local)
				schedule.Remind.Time = datetime
				h.lkDb.UpdateSchedule(schedule)
				messages = append(messages, linebot.NewTextMessage(datetime.Format(strLayout)+"\n"+okText))
			} else {
				log.Print(err)
				messages = append(messages, linebot.NewTextMessage(errorText))
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
	action := linebot.NewDatetimePickerAction(
		pickerLabel,
		postback,
		"datetime",
		now.Format(lineLayout),
		max.Format(lineLayout),
		now.Format(lineLayout),
	)
	template := linebot.NewButtonsTemplate("", "", pickerText, action)

	return linebot.NewTemplateMessage(pickerLabel, template)
}
