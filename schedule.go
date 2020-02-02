package main

import (
	"context"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Schedule struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	UserID  string             `bson:"user_id"`
	Content string             `bson:"content"`
}

func newScedule(userID string, content string) *Schedule {
	return &Schedule{
		UserID:  userID,
		Content: strings.TrimSpace(content),
	}
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
