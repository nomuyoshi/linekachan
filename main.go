package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/linebot"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	defer mongoClient.Disconnect(context.Background())
	db = mongoClient.Database("linekachan")

	bot, err := linebot.New(os.Getenv("LINE_CHANNEL_SECRET"), os.Getenv("LINE_ACCESS_TOKEN"))
	if err != nil {
		log.Fatal("Error linebot new:", err)
	}

	callbackHandler := &CallbackHandler{bot: bot}
	router := mux.NewRouter().StrictSlash(true)
	router.Handle("/callback", callbackHandler).Methods("POST")
	http.Handle("/", router)

	log.Print("Web Server starting port: 3000")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
