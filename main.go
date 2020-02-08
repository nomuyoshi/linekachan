package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bamzi/jobrunner"
	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
	"github.com/line/line-bot-sdk-go/linebot"

	_ "github.com/lib/pq"
	"gopkg.in/gorp.v2"
)

func init() {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	time.Local = jst
}

func main() {
	db, _ := sql.Open("postgres", os.Getenv("POSTGRESQL_DSN"))
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}
	lkDb := NewLineKachanDb(dbmap)

	if err := lkDb.CreateTables(); err != nil {
		log.Fatal("Error create table:", err)
	}

	bot, err := linebot.New(os.Getenv("LINE_CHANNEL_SECRET"), os.Getenv("LINE_ACCESS_TOKEN"))
	if err != nil {
		log.Fatal("Error linebot new:", err)
	}

	jobrunner.Start()
	jobrunner.Schedule("@every 5m", Reminder{})

	callbackHandler := &CallbackHandler{bot: bot, lkDb: lkDb}
	router := mux.NewRouter().StrictSlash(true)
	router.Handle("/callback", callbackHandler).Methods("POST")
	http.Handle("/", router)

	log.Print("Web Server starting port: 3000")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
