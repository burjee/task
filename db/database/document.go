package database

import (
	"time"

	pb "task/grpc"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Task struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	Title     string        `bson:"title,omitempty"`
	Status    pb.Status     `bson:"status"`
	CreatedAt string        `bson:"created_at,omitempty"`
	UpdatedAt string        `bson:"updated_at,omitempty"`
}

func NewTask(title string) *Task {
	now := time.Now().UTC().Format(time.DateTime)
	return &Task{
		Title:     title,
		Status:    pb.Status_PENDING,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
