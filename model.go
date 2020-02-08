package main

import (
	"log"
	"strconv"
	"strings"
	"time"

	"gopkg.in/gorp.v2"
)

// LineKachanDb はDatabase操作をする型
type LineKachanDb struct {
	dbmap *gorp.DbMap
}

// NewLineKachanDb は新しいLineKachanDbを作成する
func NewLineKachanDb(dbmap *gorp.DbMap) *LineKachanDb {
	return &LineKachanDb{dbmap: dbmap}
}

// CreateTables はテーブル作成をする
func (lkDb *LineKachanDb) CreateTables() error {
	lkDb.dbmap.AddTableWithName(Schedule{}, "schedules").SetKeys(true, "id")
	return lkDb.dbmap.CreateTablesIfNotExists()
}

// FindScheduleBy はScheduleをid, user_idをもとに取得する
func (lkDb *LineKachanDb) FindScheduleBy(id int64, userID string) (*Schedule, error) {
	var schedule Schedule
	err := lkDb.dbmap.SelectOne(&schedule, "select * from schedules where id=$1 and user_id=$2", id, userID)
	if err != nil {
		return &schedule, err
	}
	return &schedule, nil
}

// AddSchedule はデータベースにScheduleを追加する
func (lkDb *LineKachanDb) AddSchedule(schedule *Schedule) error {
	return lkDb.dbmap.Insert(schedule)
}

// UpdateSchedule はデータベースのScheduleを更新する
func (lkDb *LineKachanDb) UpdateSchedule(schedule *Schedule) (int64, error) {
	log.Print(schedule)
	return lkDb.dbmap.Update(schedule)
}

// Schedule はリマインドスケジュールを管理する
type Schedule struct {
	ID      int64     `db:"id, primarykey, autoincrement"`
	UserID  string    `db:"user_id, notnull"`
	Content string    `db:"content, notnull"`
	Remind  time.Time `db:"remind"`
	Created time.Time `db:"created_at, notnull"`
}

// PreInsert はDBへのInsert前のフック
func (s *Schedule) PreInsert(sql gorp.SqlExecutor) error {
	s.Created = time.Now()
	return nil
}

// NewSchedule は新しいリマインドスケジュールを作成する
func NewSchedule(userID string, content string) *Schedule {
	return &Schedule{
		UserID:  userID,
		Content: strings.TrimSpace(content),
	}
}

func (s *Schedule) postbackData() string {
	return "scheduleId=" + strconv.FormatInt(s.ID, 10)
}
