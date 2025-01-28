package server

import (
	"context"
	"errors"
	"regexp"
	"task/db/database"
	pb "task/grpc"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

var valid_id = regexp.MustCompile(`^[a-z0-9]{24}$`)
var valid_title = regexp.MustCompile(`^[a-zA-Z0-9]{1,10}$`)

type server struct {
	pb.UnimplementedTaskManagerServer
	Repository Repository
}

func New(coll *mongo.Collection) *server {
	return &server{Repository: &TaskRepository{coll}}
}

func (s *server) List(ctx context.Context, _ *pb.EmptyRequest) (*pb.ListReply, error) {
	db_tasks, err := s.Repository.Find100(ctx)
	if err != nil {
		return nil, err
	}

	grpc_tasks := convertTasks(db_tasks)
	return &pb.ListReply{Tasks: grpc_tasks}, nil
}

func (s *server) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddReply, error) {
	if !valid_title.MatchString(req.GetTitle()) {
		return nil, errors.New("input")
	}

	_, err := s.Repository.InsertOne(ctx, req.GetTitle())
	if err != nil {
		return nil, err
	}

	return &pb.AddReply{Ok: true}, nil
}

func (s *server) Update(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateReply, error) {
	if !valid_id.MatchString(req.GetId()) || !(req.GetStatus() == pb.Status_PENDING || req.GetStatus() == pb.Status_IN_PROGRESS || req.GetStatus() == pb.Status_COMPLETED) {
		return nil, errors.New("input")
	}

	result, err := s.Repository.UpdateByID(ctx, req.GetId(), req.GetStatus())
	if err != nil {
		return nil, err
	}

	if result.ModifiedCount == 0 {
		return nil, errors.New("no data")
	}

	return &pb.UpdateReply{Ok: true}, nil
}

func (s *server) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteReply, error) {
	if !valid_id.MatchString(req.GetId()) {
		return nil, errors.New("input")
	}

	result, err := s.Repository.DeleteOne(ctx, req.GetId())
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
