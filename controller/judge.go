package controller

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/3FanYu/Judges321-backend/database"
	"github.com/3FanYu/Judges321-backend/model"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var wg = sync.WaitGroup{}

func CreateJudge(c *gin.Context) {
	// 接body的資料
	var judges model.Judge
	err := c.Bind(&judges)
	if err != nil {
		log.Fatal(err)
	}
	//資料放進[]interfce{}
	var tmp []interface{}
	for _, v := range judges.Names {
		var judge = model.Judge{
			EventID: judges.EventID,
			Name:    v,
			RowNum:  judges.RowNum,
		}
		tmp = append(tmp, judge)
	}
	// 開存
	collection := database.Db.Collection("judges")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	res, err := collection.InsertMany(ctx, tmp)
	if err != nil {
		res = nil
	}
	c.JSON(200, gin.H{
		"message": true,
		"judgeID": res.InsertedIDs,
	})
}

func GetJudge(c *gin.Context) {
	// 接body的資料
	id := c.Query("judgeID")
	judgeID, err := primitive.ObjectIDFromHex(id) // 參數轉成 objectID
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(judgeID)
	//取得judge資料
	collection := database.Db.Collection("judges")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	res := collection.FindOne(ctx, bson.M{"_id": judgeID})
	// 將查詢結果放入struct
	var judge model.Judge
	res.Decode(&judge)

	eventID, err := primitive.ObjectIDFromHex(judge.EventID)
	if err != nil {
		log.Fatal(err)
	}
	collection = database.Db.Collection("events")
	res = collection.FindOne(ctx, bson.M{"_id": eventID}) // 取得活動資料
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"number", 1}, {"row", 1}})
	collection = database.Db.Collection("scores")
	cursor, err := collection.Find(ctx, bson.M{"judgeID": id}, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	var scores []model.Score
	if err = cursor.All(ctx, &scores); err != nil {
		log.Fatal(err)
	}
	var arrangedScores []interface{}
	go arrangeScores(scores, &arrangedScores, judge.RowNum)
	wg.Add(1)

	// tmp:= make(map[int][]model.Score)
	// tmp := [][]model.Score{}
	// tmp := [][]model.Score{}
	// for _, s := range scores {
	// 		tmp[s.Number] = append(tmp[s.Number], s)
	// }
	// fmt.Println(tmp)

	var event model.Event
	res.Decode(&event)
	var judgeInfo = model.JudgeInfo{
		EventName:  event.Name,
		EventOwner: event.Owner,
		JudgeName:  judge.Name,
		RowNum:     judge.RowNum,
	}
	wg.Wait()
	c.JSON(200, gin.H{
		"message":   true,
		"judgeInfo": judgeInfo,
		"scores":    arrangedScores,
	})
}

// 按照 row、 number 重新排序所有分數，中間有空的分數直接補入空值
func arrangeScores(scores []model.Score, arrangedScores *[]interface{}, rowNum int) {
	var scoreArray *[]interface{} = arrangedScores
	var subArray []interface{}
	r, n := 1, 1

	// for i:=0;i<greatestNumber;i++{
	// 	var subArray []interface{}
	// 	for y:=0;y<rowNum;y++{
	// 		if scores[i].Number
	// 	}
	// }

	for i := 0; i < len(scores); {
		if scores[i].Number == n && scores[i].Row == r { //該號碼該排有分數就插入
			subArray = append(subArray, scores[i])
			fmt.Println("append 1, currentRow: ", r, " currentNum: ", n)
			r++
			i++
		} else { //該號碼該排沒分數就插入nil
			subArray = append(subArray, nil)
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
					subArray = append(subArray, nil)
					r++
				}
				*scoreArray = append(*scoreArray, subArray)
			}
		}
	}
	wg.Done()
}
