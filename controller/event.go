package controller

import (
	"context"
	"log"
	"time"

	"github.com/3FanYu/Judges321-backend/database"
	"github.com/3FanYu/Judges321-backend/model"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateEvent(c *gin.Context) {
	// 接body的資料
	var event model.Event
	err := c.BindJSON(&event)
	if err != nil {
		log.Fatal(err)
	}
	// 開存
	collection := database.Db.Collection("events")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	res, err := collection.InsertOne(ctx, event)
	if err != nil {
		res = nil
	}
	c.JSON(200, gin.H{
		"message": true,
		"eventID":res.InsertedID.(primitive.ObjectID).Hex(),
	})
}
