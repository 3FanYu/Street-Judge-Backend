package controller

import (
	"context"
	"log"
	"time"

	"github.com/3FanYu/Judges321-backend/database"
	"github.com/3FanYu/Judges321-backend/model"
	"github.com/gin-gonic/gin"
)

func AddScore(c *gin.Context) {
	var newScore model.Score
	err := c.BindJSON(&newScore)
	if err != nil {
		log.Fatal(err)
	}
	collection := database.Db.Collection("scores")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	res, err := collection.InsertOne(ctx, newScore)
	if err != nil {
		res = nil
	}
	c.JSON(200,gin.H{
		"message":res.InsertedID,
	})

}
