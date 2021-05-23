package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/3FanYu/Judges321-backend/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	MongoClient *mongo.Client
	Db          *mongo.Database
)

func ConnectDB(conf config.MongoConfiguration) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(conf.Server))
	if err != nil {
		panic(err)
	} else {
		fmt.Println("connected to db successfully")

	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	MongoClient = client
	Db = client.Database(conf.Database)
}
