package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
	"gopkg.in/gorp.v2"
)

// リマインド登録に必要なキーワード
const (
	lineLayout = "2006-01-02T15:04"
	strLayout  = "2006/01/02 15:04"
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

				if strings.HasPrefix(resText, Prefix) {
					schedule := NewSchedule(event.Source.UserID, strings.Replace(resText, Prefix, "", 1))
					if err := h.lkDb.AddSchedule(schedule); err != nil {
						break
					}

					messages = append(messages, buildDatetimePickerMessage(schedule.PostbackData()))
				} else {
					messages = append(messages, linebot.NewTextMessage(Usage))
				}
			} else {
				messages = append(messages, linebot.NewTextMessage(Usage))
			}
		case linebot.EventTypePostback:
			if !strings.HasPrefix(event.Postback.Data, "scheduleId=") {
				continue
			}
			scheduleID, _ := strconv.ParseInt(strings.TrimPrefix(event.Postback.Data, "scheduleId="), 10, 64)
			var schedule Schedule
			if err := h.lkDb.SelectOneScheduleBy(&schedule, scheduleID, event.Source.UserID); err == nil {
				datetime, _ := time.ParseInLocation(lineLayout, event.Postback.Params.Datetime, time.Local)
				schedule.Remind = gorp.NullTime{Time: datetime, Valid: true}
				h.lkDb.UpdateSchedule(&schedule)
				messages = append(messages, linebot.NewTextMessage(datetime.Format(strLayout)+"\n"+OkText))
			} else {
				log.Print(err)
				messages = append(messages, linebot.NewTextMessage(ErrorText))
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
		PickerLabel,
		postback,
		"datetime",
		now.Format(lineLayout),
		max.Format(lineLayout),
		now.Format(lineLayout),
	)
	template := linebot.NewButtonsTemplate("", "", PickerText, action)

	return linebot.NewTemplateMessage(PickerLabel, template)
}
