package main

import (
	"context"
	"time"

	"github.com/3FanYu/Judges321-backend/config"
	"github.com/3FanYu/Judges321-backend/controller"
	"github.com/3FanYu/Judges321-backend/database"
	"github.com/gin-gonic/gin"
)

func main() {
	// set up database
	conf := config.GetConfig()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	database.ConnectDB(conf.Mongo)
	defer func() {
		if err := database.MongoClient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	router := gin.Default()
	router.POST("/api/event", controller.CreateEvent)
	router.POST("/api/score", controller.AddScore)
	router.Run("localhost:8080")
}
