package server

import (
	"context"
	"task/db/database"
	pb "task/grpc"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var object_id = bson.NewObjectID()
var err_input = "input"

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Find100(ctx context.Context) ([]database.Task, error) {
	args := m.Called(ctx)
	return args.Get(0).([]database.Task), args.Error(1)
}

func (m *MockRepository) InsertOne(ctx context.Context, title string) (*mongo.InsertOneResult, error) {
	args := m.Called(ctx, title)
	return args.Get(0).(*mongo.InsertOneResult), args.Error(1)
}

func (m *MockRepository) UpdateByID(ctx context.Context, id string, status pb.Status) (*mongo.UpdateResult, error) {
	args := m.Called(ctx, id, status)
	return args.Get(0).(*mongo.UpdateResult), args.Error(1)
}

func (m *MockRepository) DeleteOne(ctx context.Context, id string) (*mongo.DeleteResult, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*mongo.DeleteResult), args.Error(1)
}

func TestList(t *testing.T) {
	db_tasks := []database.Task{{ID: object_id, Title: "Title", Status: 1, CreatedAt: "2025-01-01 00:00:00", UpdatedAt: "2025-01-01 12:00:00"}}
	pb_tasks := []*pb.Task{{Id: object_id.Hex(), Title: "Title", Status: 1, CreatedAt: "2025-01-01 00:00:00", UpdatedAt: "2025-01-01 12:00:00"}}

	mock_repository := new(MockRepository)
	mock_repository.On("Find100", mock.Anything).Return(db_tasks, nil)
	server := &server{Repository: mock_repository}

	server_reply, err := server.List(context.Background(), &pb.EmptyRequest{})
	expected_reply := &pb.ListReply{Tasks: pb_tasks}

	assert.Equal(t, expected_reply.String(), server_reply.String())
	assert.Nil(t, err)
	mock_repository.AssertExpectations(t)
}

func TestAdd(t *testing.T) {
	mock_repository := new(MockRepository)
	mock_repository.On("InsertOne", mock.Anything, mock.AnythingOfType("string")).Return(&mongo.InsertOneResult{}, nil)
	server := &server{Repository: mock_repository}

	req := &pb.AddRequest{Title: "Title"}
	server_reply, err := server.Add(context.Background(), req)

	assert.True(t, server_reply.GetOk())
	assert.Nil(t, err)
	mock_repository.AssertExpectations(t)
}

func TestAdd_InvalidInput(t *testing.T) {
	mock_repository := new(MockRepository)
	server := &server{Repository: mock_repository}

	server_reply, err := server.Add(context.Background(), nil)

	assert.Nil(t, server_reply)
	assert.EqualError(t, err, err_input)
	mock_repository.AssertExpectations(t)
}

func TestUpdate(t *testing.T) {
	mock_repository := new(MockRepository)
	mock_repository.On("UpdateByID", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("grpc.Status")).Return(&mongo.UpdateResult{ModifiedCount: 1}, nil)
	server := &server{Repository: mock_repository}

	req := &pb.UpdateRequest{Id: object_id.Hex(), Status: pb.Status_COMPLETED}
	server_reply, err := server.Update(context.Background(), req)

	assert.True(t, server_reply.GetOk())
	assert.Nil(t, err)
	mock_repository.AssertExpectations(t)
}

func TestUpdate_InvalidInput_ID(t *testing.T) {
	mock_repository := new(MockRepository)
	server := &server{Repository: mock_repository}

	req := &pb.UpdateRequest{Id: "Error ID", Status: pb.Status_COMPLETED}
	server_reply, err := server.Update(context.Background(), req)

	assert.Nil(t, server_reply)
	assert.EqualError(t, err, err_input)
	mock_repository.AssertExpectations(t)
}

func TestUpdate_InvalidInput_Status(t *testing.T) {
	mock_repository := new(MockRepository)
	server := &server{Repository: mock_repository}

	req := &pb.UpdateRequest{Id: object_id.Hex(), Status: -1}
	server_reply, err := server.Update(context.Background(), req)

	assert.Nil(t, server_reply)
	assert.EqualError(t, err, err_input)
	mock_repository.AssertExpectations(t)
}

func TestDelete(t *testing.T) {
	mock_repository := new(MockRepository)
	mock_repository.On("DeleteOne", mock.Anything, mock.AnythingOfType("string")).Return(&mongo.DeleteResult{DeletedCount: 1}, nil)
	server := &server{Repository: mock_repository}

	req := &pb.DeleteRequest{Id: object_id.Hex()}
	server_reply, err := server.Delete(context.Background(), req)

	assert.True(t, server_reply.GetOk())
	assert.Nil(t, err)
	mock_repository.AssertExpectations(t)

}

func TestDelete_InvalidInput(t *testing.T) {
	mock_repository := new(MockRepository)
	server := &server{Repository: mock_repository}

	req := &pb.DeleteRequest{Id: "Error ID"}
	server_reply, err := server.Delete(context.Background(), req)

	assert.Nil(t, server_reply)
	assert.EqualError(t, err, err_input)
	mock_repository.AssertExpectations(t)
}
