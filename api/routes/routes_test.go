package routes

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	pb "task/grpc"

	"context"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

var msg_ok = `{"message":"ok"}`
var err_id = `{"error":"id error"}`
var err_title = `{"error":"title error"}`
var err_status = `{"error":"status error"}`

type MockGRPC struct {
	mock.Mock
}

func (m *MockGRPC) List(ctx context.Context, in *pb.EmptyRequest, opts ...grpc.CallOption) (*pb.ListReply, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*pb.ListReply), args.Error(1)
}

func (m *MockGRPC) Add(ctx context.Context, in *pb.AddRequest, opts ...grpc.CallOption) (*pb.AddReply, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*pb.AddReply), args.Error(1)
}

func (m *MockGRPC) Update(ctx context.Context, in *pb.UpdateRequest, opts ...grpc.CallOption) (*pb.UpdateReply, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*pb.UpdateReply), args.Error(1)
}

func (m *MockGRPC) Delete(ctx context.Context, in *pb.DeleteRequest, opts ...grpc.CallOption) (*pb.DeleteReply, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*pb.DeleteReply), args.Error(1)
}

func init() {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard
	gin.DefaultErrorWriter = io.Discard
}

func TestList(t *testing.T) {
	task := &pb.Task{Id: "507f1f77bcf86cd799439011", Title: "title", Status: pb.Status_COMPLETED, CreatedAt: "2025-01-01 00:00:00", UpdatedAt: "2025-01-01 12:00:00"}
	tasks := []*pb.Task{task}

	g := gin.Default()
	mock_grpc := new(MockGRPC)
	mock_grpc.On("List", mock.Anything, mock.Anything, mock.Anything).Return(&pb.ListReply{Tasks: tasks}, nil)
	Setup(g, mock_grpc)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/task", nil)
	g.ServeHTTP(w, req)

	task_json, _ := protojson.MarshalOptions{UseEnumNumbers: true, EmitDefaultValues: true}.Marshal(&pb.ListReply{Tasks: tasks})

	assert.Equal(t, 200, w.Code)
	assert.JSONEq(t, string(task_json), w.Body.String())
	mock_grpc.AssertExpectations(t)
}

func TestAdd(t *testing.T) {
	g := gin.Default()
	mock_grpc := new(MockGRPC)
	mock_grpc.On("Add", mock.Anything, mock.Anything, mock.Anything).Return(&pb.AddReply{Ok: true}, nil)
	Setup(g, mock_grpc)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/task", strings.NewReader(`{"title":"title"}`))
	g.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.JSONEq(t, msg_ok, w.Body.String())
	mock_grpc.AssertExpectations(t)
}

func TestAdd_InvalidInput(t *testing.T) {
	g := gin.Default()
	mock_grpc := new(MockGRPC)
	Setup(g, mock_grpc)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/task", nil)
	g.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	assert.JSONEq(t, err_title, w.Body.String())
	mock_grpc.AssertExpectations(t)
}

func TestUpdate(t *testing.T) {
	g := gin.Default()
	mock_grpc := new(MockGRPC)
	mock_grpc.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(&pb.UpdateReply{Ok: true}, nil)
	Setup(g, mock_grpc)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/task/507f1f77bcf86cd799439011", strings.NewReader(`{"status":1}`))
	g.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.JSONEq(t, msg_ok, w.Body.String())
	mock_grpc.AssertExpectations(t)
}

func TestUpdate_InvalidInput_ID(t *testing.T) {
	g := gin.Default()
	mock_grpc := new(MockGRPC)
	Setup(g, mock_grpc)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/task/ERROR_ID", strings.NewReader(`{"status":1}`))
	g.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	assert.JSONEq(t, err_id, w.Body.String())
	mock_grpc.AssertExpectations(t)
}

func TestUpdate_InvalidInput_Status(t *testing.T) {
	g := gin.Default()
	mock_grpc := new(MockGRPC)
	Setup(g, mock_grpc)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/task/507f1f77bcf86cd799439011", nil)
	g.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	assert.JSONEq(t, err_status, w.Body.String())
	mock_grpc.AssertExpectations(t)
}

func TestDelete(t *testing.T) {
	g := gin.Default()
	mock_grpc := new(MockGRPC)
	mock_grpc.On("Delete", mock.Anything, mock.Anything, mock.Anything).Return(&pb.DeleteReply{Ok: true}, nil)
	Setup(g, mock_grpc)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/task/507f1f77bcf86cd799439011", nil)
	g.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.JSONEq(t, msg_ok, w.Body.String())
	mock_grpc.AssertExpectations(t)
}

func TestDelete_InvalidInput(t *testing.T) {
	g := gin.Default()
	mock_grpc := new(MockGRPC)
	Setup(g, mock_grpc)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/task/ERROR_ID", nil)
	g.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	assert.JSONEq(t, err_id, w.Body.String())
	mock_grpc.AssertExpectations(t)
}
