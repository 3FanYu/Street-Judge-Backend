package controller

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/3FanYu/Judges321-backend/database"
	"github.com/3FanYu/Judges321-backend/model"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SettleScore(c *gin.Context) {
	eventID := c.Query("eventID")
	judges := getAllJudges(&eventID)
	ch := make(chan model.JudgeScore)
	for i := 0; i < len(*judges); i++ {
		go getAllScores(&((*judges)[i]), ch)
	}
	var result model.Settlement
	for x := 0; x < len(*judges); x++ {
		result.JudgeScore = append(result.JudgeScore, <-ch)
	}
	c.JSON(200, gin.H{
		"message":    true,
		"settlement": result,
	})
}

func getAllJudges(eventID *string) *[]model.Judge {
	collection := database.Db.Collection("judges")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	findOptions := options.Find().SetSort(bson.D{{"name", 1}})
	cursor, err := collection.Find(ctx, bson.M{"eventID": eventID}, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	var judges []model.Judge
	if err = cursor.All(ctx, &judges); err != nil {
		log.Fatal(err)
	}
	fmt.Println(judges)
	return &judges
}

func getAllScores(judge *model.Judge, ch chan model.JudgeScore) {
	judgeID := judge.ID.Hex()
	collection := database.Db.Collection("scores")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	findOptions := options.Find().SetSort(bson.D{{"number", 1}, {"row", 1}})
	cursor, err := collection.Find(ctx, bson.M{"judgeID": judgeID}, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	var scores []model.Score
	if err = cursor.All(ctx, &scores); err != nil {
		log.Fatal(err)
	}
	if len(scores) > 0 {
		var judgeScore = model.JudgeScore{
			ID:     judge.ID,
			Name:   judge.Name,
			RowNum: judge.RowNum,
			Scores: scores,
		}
		ch <- judgeScore
	} else {
		var judgeScore = model.JudgeScore{
			ID:     judge.ID,
			Name:   judge.Name,
			RowNum: judge.RowNum,
			Scores: []model.Score{},
		}
		ch <- judgeScore
	}

}
