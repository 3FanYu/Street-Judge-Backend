package main

import (
	"context"
	"time"

	"github.com/3FanYu/Judges321-backend/config"
	"github.com/3FanYu/Judges321-backend/controller"
	"github.com/3FanYu/Judges321-backend/database"
	"github.com/gin-contrib/cors"
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
	router.Use(cors.Default())
	router.POST("/api/event", controller.CreateEvent)
	router.GET("/api/event", controller.GetEvent)
	router.POST("/api/score", controller.AddScore)
	router.POST("/api/judge", controller.CreateJudge)
	router.GET("/api/judge", controller.GetJudge)
	router.GET("/api/settlement", controller.SettleScore)

	router.Run("localhost:8080")
}
