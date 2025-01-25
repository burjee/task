package main

import (
	"log"
	"net"
	_ "task/db/config"
	"task/db/database"
	"task/db/server"
	pb "task/grpc"

	"google.golang.org/grpc"
)

func main() {
	db, close := database.New()
	defer close()

	grpc_server := grpc.NewServer()
	task_server := server.New(db)
	pb.RegisterTaskManagerServer(grpc_server, task_server)

	listen, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	log.Printf("server listening at %v", listen.Addr())
	log.Fatal(grpc_server.Serve(listen))
}
