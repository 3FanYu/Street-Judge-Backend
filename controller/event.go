package controller

import (
	"log"

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
	res, err := event.CreateEvent()
	if err != nil {
		c.JSON(400, gin.H{
			"result":  false,
			"message": "錯誤錯誤",
		})
	}
	c.JSON(200, gin.H{
		"result":  true,
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
	var event model.Event
	event.GetEvent(&eventID)
	judges := getAllJudges(&ID)
	event.Judges = *judges
	c.JSON(200, gin.H{
		"result": true,
		"event":  &event,
	})
}
