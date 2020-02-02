package main

import (
	"context"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Schedule struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	UserID         string             `bson:"user_id"`
	Content        string             `bson:"content"`
	RemindDatetime time.Time          `bson:"remind_datetime"`
}

func newSchedule(userID string, content string) *Schedule {
	return &Schedule{
		UserID:  userID,
		Content: strings.TrimSpace(content),
	}
}

func findScheduleBy(id string, userID string) (*Schedule, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	var schedule Schedule
	objectID, _ := primitive.ObjectIDFromHex(id)
	err := db.Collection("schedules").FindOne(ctx, bson.M{"_id": objectID, "user_id": userID}).Decode(&schedule)
	if err != nil {
		return nil, err
	}
	return &schedule, nil
}

func (s *Schedule) update(datetime time.Time) error {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	filter := bson.M{"_id": s.ID}
	update := bson.M{"$set": bson.M{"remind_datetime": datetime}}

	_, err := db.Collection("schedules").UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	s.RemindDatetime = datetime
	return nil
}

func (s *Schedule) create() error {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	res, err := db.Collection("schedules").InsertOne(ctx, &s)
	if err != nil {
		return err
	}

	s.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

func (s *Schedule) postbackData() string {
	return "scheduleId=" + s.ID.Hex()
}
