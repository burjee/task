package server

import (
	"context"
	"errors"
	"regexp"
	"task/db/database"
	pb "task/grpc"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var valid_id = regexp.MustCompile(`^[a-z0-9]{24}$`)
var valid_title = regexp.MustCompile(`^[a-zA-Z0-9]{1,10}$`)

type server struct {
	pb.UnimplementedTaskManagerServer
	db   *mongo.Client
	coll *mongo.Collection
}

func New(db *mongo.Client) *server {
	coll := db.Database("task").Collection("tasks")
	return &server{db: db, coll: coll}
}

func (s *server) List(_ context.Context, _ *pb.EmptyRequest) (*pb.ListReply, error) {
	cur, err := s.coll.Find(context.Background(), bson.D{}, options.Find().SetLimit(100))
	if err != nil {
		return nil, err
	}

	var db_tasks []database.Task
	// This method will close the cursor after retrieving all documents. no need cur.Close()
	if err := cur.All(context.Background(), &db_tasks); err != nil {
		return nil, err
	}

	grpc_tasks := convertTasks(db_tasks)
	return &pb.ListReply{Tasks: grpc_tasks}, nil
}

func (s *server) Add(_ context.Context, req *pb.AddRequest) (*pb.AddReply, error) {
	if !valid_title.MatchString(req.GetTitle()) {
		return nil, errors.New("input")
	}

	doc := database.NewTask(req.GetTitle())
	_, err := s.coll.InsertOne(context.Background(), doc)
	if err != nil {
		return nil, err
	}

	return &pb.AddReply{Ok: true}, nil
}

func (s *server) Update(_ context.Context, req *pb.UpdateRequest) (*pb.UpdateReply, error) {
	if !valid_id.MatchString(req.GetId()) || !(req.GetStatus() == pb.Status_PENDING || req.GetStatus() == pb.Status_IN_PROGRESS || req.GetStatus() == pb.Status_COMPLETED) {
		return nil, errors.New("input")
	}

	oid, err := bson.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, err
	}

	update := bson.D{{Key: "$set", Value: bson.D{{Key: "status", Value: req.GetStatus()}, {Key: "updated_at", Value: time.Now().UTC().Format(time.DateTime)}}}}
	result, err := s.coll.UpdateByID(context.Background(), oid, update)
	if err != nil {
		return nil, err
	}

	if result.ModifiedCount == 0 {
		return nil, errors.New("no data")
	}

	return &pb.UpdateReply{Ok: true}, nil
}

func (s *server) Delete(_ context.Context, req *pb.DeleteRequest) (*pb.DeleteReply, error) {
	if !valid_id.MatchString(req.GetId()) {
		return nil, errors.New("input")
	}

	oid, err := bson.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, err
	}

	filter := bson.D{{Key: "_id", Value: oid}}
	result, err := s.coll.DeleteOne(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	if result.DeletedCount == 0 {
		return nil, errors.New("no data")
	}

	return &pb.DeleteReply{Ok: true}, nil
}

func convertTasks(db_tasks []database.Task) []*pb.Task {
	grpc_tasks := make([]*pb.Task, 0, len(db_tasks))
	for _, db_task := range db_tasks {
		grpc_tasks = append(grpc_tasks, convertSingleTask(db_task))
	}

	return grpc_tasks
}

func convertSingleTask(db_task database.Task) *pb.Task {
	return &pb.Task{
		Id:        db_task.ID.Hex(),
		Title:     db_task.Title,
		Status:    db_task.Status,
		CreatedAt: db_task.CreatedAt,
		UpdatedAt: db_task.UpdatedAt,
	}
}
