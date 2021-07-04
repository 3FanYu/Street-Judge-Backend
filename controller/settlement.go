package controller

import (
	"context"
	"log"
	"time"

	"github.com/3FanYu/Judges321-backend/database"
	"github.com/3FanYu/Judges321-backend/model"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SettleScore(c *gin.Context) {
	eventID := c.Query("eventID")     // 取得 eventID
	judges := getAllJudges(&eventID)  // 取得此 eventID 的所有judge
	ch := make(chan model.JudgeScore) // 創建channel
	for i := 0; i < len(*judges); i++ {
		go getAllScores(&((*judges)[i]), ch) //同步取得多位judge的所有分數，並使用channel回傳。
	}
	result := make([]interface{}, 0)    // 創造空slice並依序把channel收到的回傳值插入
	for x := 0; x < len(*judges); x++ { // 有幾位judge就接收channel幾次
		result = append(result, <-ch)
	}
	c.JSON(200, gin.H{ //回傳
		"message":    true,
		"settlement": result,
	})
}

// 取得此 eventID 的所有judge
func getAllJudges(eventID *string) *[]model.Judge {
	collection := database.Db.Collection("judges")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	findOptions := options.Find().SetSort(bson.D{{"name", 1}})
	// 查詢所有eventID符合的judge
	cursor, err := collection.Find(ctx, bson.M{"eventID": eventID}, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	var judges []model.Judge
	if err = cursor.All(ctx, &judges); err != nil {
		log.Fatal(err)
	}
	return &judges
}

// 取得該judge的所有分數
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
	// 如果沒有分數就回傳空陣列
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
