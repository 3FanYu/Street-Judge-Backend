package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/3FanYu/Judges321-backend/config"
	"github.com/3FanYu/Judges321-backend/controller"
	"github.com/3FanYu/Judges321-backend/database"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// load .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
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

	port := os.Getenv("PORT")
	router := gin.Default()
	// router.Use(cors.New(config.CorsConfig()))
	router.Use(cors.Default())
	router.GET("/api", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.POST("/api/event", controller.CreateEvent)
	router.GET("/api/event", controller.GetEvent)
	router.POST("/api/score", controller.AddScore)
	router.POST("/api/judge", controller.CreateJudge)
	router.GET("/api/judge", controller.GetJudge)
	router.GET("/api/settlement", controller.SettleScore)

	router.Run("localhost:" + port)
}
