package controller

import (
	"context"
	"log"
	"time"

	"github.com/3FanYu/Judges321-backend/database"
	"github.com/3FanYu/Judges321-backend/model"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func AddScore(c *gin.Context) {
	var newScore model.Score
	err := c.ShouldBindJSON(&newScore)
	if err != nil {
		c.JSON(400, gin.H{
			"result":  false,
			"message": err.Error(),
		})
		return
	}
	isExisted, existedScore := checkIfExist(newScore)
	if isExisted {
		updateScore(&newScore, existedScore, c)
	} else {
		insertScore(&newScore, c)
	}
}

func checkIfExist(score model.Score) (bool, *model.Score) {
	judgeID := score.JudgeID
	row := score.Row
	number := score.Number
	var existedScore model.Score
	collection := database.Db.Collection("scores")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	filter := bson.M{"judgeID": judgeID, "row": row, "number": number}
	err := collection.FindOne(ctx, filter).Decode(&existedScore)
	if err != nil {
		return false, &existedScore
	}
	return true, &existedScore
}

func updateScore(score *model.Score, existedScore *model.Score, c *gin.Context) {
	existedScoreID := (*existedScore).ID
	point := (*score).Point
	collection := database.Db.Collection("scores")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	res, err := collection.UpdateOne(ctx, bson.M{"_id": existedScoreID}, bson.D{{"$set", bson.D{{"point", point}}}})
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(200, gin.H{
		"result":     true,
		"insertedID": res,
	})
}

func insertScore(newScore *model.Score, c *gin.Context) {
	collection := database.Db.Collection("scores")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	res, err := collection.InsertOne(ctx, *newScore)
	if err != nil {
		res = nil
	}
	c.JSON(200, gin.H{
		"result":     true,
		"insertedID": res.InsertedID,
	})
}
