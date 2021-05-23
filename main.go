package main

import (
	"context"
	"time"

	"github.com/3FanYu/Judges321-backend/config"
	"github.com/3FanYu/Judges321-backend/database"
)

func main() {
	conf := config.GetConfig()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	database.ConnectDB(conf.Mongo)
	defer func() {
		if err := database.MongoClient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}