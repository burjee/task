package libs

import (
	pb "task/grpc"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewTaskClient() (*grpc.ClientConn, pb.TaskManagerClient) {
	host := viper.GetString("grpc.host")
	port := viper.GetString("grpc.port")
	uri := host + ":" + port

	conn, err := grpc.NewClient(uri, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic("did not connect: " + err.Error())
	}

	client := pb.NewTaskManagerClient(conn)
	return conn, client
}
