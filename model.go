package main

import (
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

// SelectOneScheduleBy はScheduleをid, user_idをもとに取得する
func (lkDb *LineKachanDb) SelectOneScheduleBy(schedule *Schedule, id int64, userID string) error {
	return lkDb.dbmap.SelectOne(schedule, "select * from schedules where id=$1 and user_id=$2", id, userID)
}

// SelectSchedulesBy は引数のStatusとremindでScheduleを絞り込んだ一覧を取得する
func (lkDb *LineKachanDb) SelectSchedulesBy(schedules *[]Schedule, status ScheduleStatus, remind time.Time) error {
	_, err := lkDb.dbmap.Select(&schedules, "select * from schedules where status=$1 and remind < $2", Scheduled, remind)
	return err
}

// AddSchedule はデータベースにScheduleを追加する
func (lkDb *LineKachanDb) AddSchedule(schedule *Schedule) error {
	return lkDb.dbmap.Insert(schedule)
}

// UpdateSchedule はデータベースのScheduleを更新する
func (lkDb *LineKachanDb) UpdateSchedule(schedule *Schedule) (int64, error) {
	return lkDb.dbmap.Update(schedule)
}

// ScheduleStatus はリマインド済みかどうかのステータス
type ScheduleStatus int

const (
	Scheduled ScheduleStatus = iota
	Reminded
)

// Schedule はリマインドスケジュールを管理する
type Schedule struct {
	ID      int64          `db:"id, primarykey, autoincrement"`
	UserID  string         `db:"user_id, notnull"`
	Content string         `db:"content, notnull"`
	Status  ScheduleStatus `db:"status, notnull"`
	Remind  gorp.NullTime  `db:"remind"`
	Created time.Time      `db:"created_at, notnull"`
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
		Status:  Scheduled,
	}
}

// PostbackData はLineメッセージに含めるpostback
func (s *Schedule) PostbackData() string {
	return "scheduleId=" + strconv.FormatInt(s.ID, 10)
}
