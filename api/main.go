package main

import (
	"log"
	_ "task/api/config"
	"task/api/libs"
	"task/api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	grpc_conn, grpc_client := libs.NewTaskClient()
	defer grpc_conn.Close()

	g := gin.Default()
	routes.Setup(g, grpc_client)

	log.Println("start server http://0.0.0.0:8000")
	log.Fatal(g.Run(":8000"))
}
