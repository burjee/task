package database

import (
	"context"
	"log"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func New() (*mongo.Client, func()) {
	host := viper.GetString("mongodb.host")
	port := viper.GetString("mongodb.port")
	uri := "mongodb://" + host + ":" + port

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	if err := client.Ping(context.Background(), readpref.Primary()); err != nil {
		panic(err)
	}

	close := func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Fatalln(err)
		}
	}

	return client, close
}
