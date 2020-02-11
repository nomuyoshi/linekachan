package main

import (
	"fmt"
	"log"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
)

// Reminder は登録されたScheduleの通知を管理する
type Reminder struct {
	lkDb *LineKachanDb
	bot  *linebot.Client
}

// Run は予定時刻のScheduleを通知する
func (r Reminder) Run() {
	fmt.Printf("Every 1 min remind\n")
	var schedules []Schedule
	if err := r.lkDb.SelectSchedulesBy(&schedules, Scheduled, time.Now()); err != nil {
		log.Print("remind select schedules error: ", err)
		return
	}

	var errorIds []int64
	for _, sch := range schedules {
		msg := linebot.NewTextMessage(RemindText + "\n" + sch.Content)
		if _, err := r.bot.PushMessage(sch.UserID, msg).Do(); err == nil {
			sch.Status = Reminded
			r.lkDb.UpdateSchedule(&sch)
		} else {
			errorIds = append(errorIds, sch.ID)
		}
	}
	log.Print("remind failed schedule ids: ", errorIds)
}
