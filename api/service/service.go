package service

import (
	"context"
	pb "task/grpc"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
)

type service struct {
	grpc_client pb.TaskManagerClient
}

func New(grpc_client pb.TaskManagerClient) *service {
	return &service{grpc_client}
}

func (s *service) List(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	reply, err := s.grpc_client.List(ctx, &pb.EmptyRequest{})
	if err != nil {
		c.AbortWithStatusJSON(503, gin.H{"error": "grpc error"})
		return
	}

	b, err := protojson.MarshalOptions{UseEnumNumbers: true, EmitDefaultValues: true}.Marshal(reply)
	if err != nil {
		c.AbortWithStatusJSON(503, gin.H{"error": "server error"})
	} else {
		c.Data(200, gin.MIMEJSON, b)
	}
}

func (s *service) Add(c *gin.Context) {
	var r AddReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "title error"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	reply, err := s.grpc_client.Add(ctx, &pb.AddRequest{Title: r.Title})
	if err != nil || !reply.GetOk() {
		c.AbortWithStatusJSON(503, gin.H{"error": "grpc error"})
		return
	}

	c.JSON(200, gin.H{"message": "ok"})
}

func (s *service) Update(c *gin.Context) {
	var u Uri
	if err := c.ShouldBindUri(&u); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "id error"})
		return
	}

	var r UpdateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "status error"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	reply, err := s.grpc_client.Update(ctx, &pb.UpdateRequest{Id: u.ID, Status: pb.Status(r.Status)})
	if err != nil || !reply.GetOk() {
		c.AbortWithStatusJSON(503, gin.H{"error": "grpc error"})
		return
	}

	c.JSON(200, gin.H{"message": "ok"})
}

func (s *service) Delete(c *gin.Context) {
	var u Uri
	if err := c.ShouldBindUri(&u); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "id error"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	reply, err := s.grpc_client.Delete(ctx, &pb.DeleteRequest{Id: u.ID})
	if err != nil || !reply.GetOk() {
		c.AbortWithStatusJSON(503, gin.H{"error": "grpc error"})
		return
	}

	c.JSON(200, gin.H{"message": "ok"})
}
