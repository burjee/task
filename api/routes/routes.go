package routes

import (
	"task/api/middleware"
	"task/api/service"

	pb "task/grpc"

	"github.com/gin-gonic/gin"
)

func Setup(g *gin.Engine, grpc_client pb.TaskManagerClient) {
	service := service.New(grpc_client)

	g.Use(middleware.Cors())

	g.Static("/assets", "./web/assets")
	g.StaticFile("/", "./web/index.html")

	task := g.Group("/api/task")
	task.GET("", service.List)
	task.POST("", service.Add)
	task.PUT("/:id", service.Update)
	task.DELETE("/:id", service.Delete)
}
