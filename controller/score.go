package controller

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/3FanYu/Judges321-backend/database"
	"github.com/3FanYu/Judges321-backend/model"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	collection := database.Db.Collection("scores")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	res, err := collection.InsertOne(ctx, newScore)
	if err != nil {
		res = nil
	}
	c.JSON(200, gin.H{
		"result":     true,
		"insertedID": res.InsertedID,
	})
}

func PatchScore(c *gin.Context) {
	id := c.Query("scoreID")
	scoreID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Fatal(err)
	}
	point, _ := strconv.ParseFloat(c.Query("point"), 32)
	// var newScore model.Score
	// err := c.BindJSON(&newScore)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	collection := database.Db.Collection("scores")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	res, err := collection.UpdateOne(ctx, bson.M{"_id": scoreID}, bson.D{{"$set", bson.D{{"point", point}}}})
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(200, gin.H{
		"result":     true,
		"insertedID": res,
	})
}
