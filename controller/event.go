package controller

import (
	"context"
	"log"
	"time"

	"github.com/3FanYu/Judges321-backend/database"
	"github.com/3FanYu/Judges321-backend/model"
	"github.com/gin-gonic/gin"
)

func CreateEvent(c *gin.Context) {
	// 接body的資料
	var eventF model.Event
	err := c.BindJSON(&eventF)
	if err != nil {
		log.Fatal(err)
	}
	// 依據rowNum 決定每個評審底下要儲存幾筆score
	var scores []model.Score
	for i := 0; i < eventF.RowNum; i++ {
		var newScore = model.NewScore()
		scores = append(scores, *newScore)
	}
	// 依據judgeName數量決定加幾筆Judge
	var judges []model.Judge
	for i := 0; i < len(eventF.JudgeName); i++ {
		var newJudge = model.Judge{
			Name:   eventF.JudgeName[i],
			Scores: scores,
		}
		judges = append(judges, newJudge)
	}
	// fmt.Printf("%+v\n", judges)
	// 塞真正要存進資料庫的資料～
	var eventB = model.Event{
		Name:     eventF.Name,
		Owner:    eventF.Owner,
		Password: eventF.Password,
		Judges:   judges,
	}
	// 開存
	collection := database.Db.Collection("events")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	res, err := collection.InsertOne(ctx, eventB)
	if err != nil {
		res = nil
	}
	c.JSON(200, gin.H{
		"message": res,
	})
}
