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
	eventID := c.Query("eventID")    // 取得 eventID
	judges := getAllJudges(&eventID) // 取得此 eventID 的所有judge
	addUps := make([][]model.Score, 0)
	ch := make(chan model.JudgeScore) // 創建channel
	for i := 0; i < len(*judges); i++ {
		go getAllScores(&((*judges)[i]), ch, &addUps) //同步取得多位judge的所有分數，並使用channel回傳。
	}
	result := make([]interface{}, 0)    // 創造空slice並依序把channel收到的回傳值插入
	for x := 0; x < len(*judges); x++ { // 有幾位judge就接收channel幾次
		result = append(result, <-ch)
	}
	c.JSON(200, gin.H{ //回傳
		"message":    true,
		"settlement": result,
		"addUps":     addUps,
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
func getAllScores(judge *model.Judge, ch chan model.JudgeScore, addUps *[][]model.Score) {
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
		var arrangedScores = make([][]model.Score, 0)
		reArrangeScores(scores, &arrangedScores, judge.RowNum)
		addUpScores(addUps, arrangedScores)
		var judgeScore = model.JudgeScore{
			ID:     judge.ID,
			Name:   judge.Name,
			RowNum: judge.RowNum,
			Scores: arrangedScores,
		}
		ch <- judgeScore
	} else {
		var judgeScore = model.JudgeScore{
			ID:     judge.ID,
			Name:   judge.Name,
			RowNum: judge.RowNum,
			Scores: [][]model.Score{},
		}
		ch <- judgeScore
	}
}

// 按照 row、 number 重新排序所有分數，中間有空的分數直接補入空值
func reArrangeScores(scores []model.Score, arrangedScores *[][]model.Score, rowNum int) {
	var scoreArray *[][]model.Score = arrangedScores
	var subArray []model.Score
	r, n := 1, 1

	for i := 0; i < len(scores); {
		if scores[i].Number == n && scores[i].Row == r { //該號碼該排有分數就插入
			subArray = append(subArray, scores[i])
			fmt.Println("append 1, currentRow: ", r, " currentNum: ", n)
			r++
			i++
		} else { //該號碼該排沒分數就插入nil
			subArray = append(subArray, model.Score{})
			fmt.Println("append empty, currentRow: ", r, " currentNum: ", n)
			r++
		}
		if r > rowNum {
			r = 1
			n++
			*scoreArray = append(*scoreArray, subArray)
			subArray = nil
		}
		//最後一筆資料如果沒有填滿所有rowNum 就把其餘的補上null
		if i == len(scores) {
			fmt.Println("last")
			if r != 1 && r <= rowNum {
				for r <= rowNum {
					subArray = append(subArray, model.Score{})
					r++
				}
				*scoreArray = append(*scoreArray, subArray)
			}
		}
	}
}

func addUpScores(addUps *[][]model.Score, allScores [][]model.Score) {
	if len(*addUps) > 0 {
		for i, scores := range allScores {
			if len(*addUps) > i {
				for j, score := range scores {
					(*addUps)[i][j].Point += score.Point
				}
			} else {
				fmt.Println("is Nil")
				(*addUps) = append((*addUps), scores)
			}
		}
	} else {
		for _, scores := range allScores {
			tmpScores := make([]model.Score, 0)
			for _, score := range scores {
				tmpScores = append(tmpScores, model.Score{
					Row:    score.Row,
					Number: score.Number,
					Point:  score.Point,
				})
			}
			*addUps = append(*addUps, tmpScores)
		}
	}
}
