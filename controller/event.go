package controller

import (
	"context"
	"log"
	"time"

	"github.com/3FanYu/Judges321-backend/database"
	"github.com/3FanYu/Judges321-backend/model"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
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
		"eventID": res.InsertedID.(primitive.ObjectID).Hex(),
	})
}

func GetEvent(c *gin.Context) {
	ID := c.Query("eventID") // 取得 eventID
	eventID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		c.JSON(401, gin.H{
			"result":  false,
			"message": "invalid eventID",
		})
	}
	event := getEvent(&eventID)
	judges := getAllJudges(&ID)
	event.Judges = *judges
	c.JSON(200, gin.H{
		"result": true,
		"event":  *event,
	})
}
func getEvent(evenID *primitive.ObjectID) *model.Event {
	collection := database.Db.Collection("events")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	res := collection.FindOne(ctx, bson.M{"_id": *evenID}, options.FindOne().SetProjection(bson.M{"password": 0}))
	var event model.Event
	res.Decode(&event)
	return &event
}
