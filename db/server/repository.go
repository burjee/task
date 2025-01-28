package server

import (
	"context"
	"task/db/database"
	pb "task/grpc"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Repository interface {
	Find100(ctx context.Context) ([]database.Task, error)
	InsertOne(ctx context.Context, title string) (*mongo.InsertOneResult, error)
	UpdateByID(ctx context.Context, id string, status pb.Status) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, id string) (*mongo.DeleteResult, error)
}

type TaskRepository struct {
	Coll *mongo.Collection
}

func (t *TaskRepository) Find100(ctx context.Context) ([]database.Task, error) {
	cur, err := t.Coll.Find(ctx, bson.D{}, options.Find().SetLimit(100))
	if err != nil {
		return nil, err
	}

	var tasks []database.Task
	// This method will close the cursor after retrieving all documents. no need cur.Close()
	err = cur.All(ctx, &tasks)
	return tasks, err
}

func (t *TaskRepository) InsertOne(ctx context.Context, title string) (*mongo.InsertOneResult, error) {
	doc := database.NewTask(title)
	return t.Coll.InsertOne(context.Background(), doc)
}

func (t *TaskRepository) UpdateByID(ctx context.Context, id string, status pb.Status) (*mongo.UpdateResult, error) {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	update := bson.D{{Key: "$set", Value: bson.D{{Key: "status", Value: status}, {Key: "updated_at", Value: time.Now().UTC().Format(time.DateTime)}}}}
	return t.Coll.UpdateByID(context.Background(), oid, update)
}

func (t *TaskRepository) DeleteOne(ctx context.Context, id string) (*mongo.DeleteResult, error) {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.D{{Key: "_id", Value: oid}}
	return t.Coll.DeleteOne(context.Background(), filter)
}
