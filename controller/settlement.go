package controller

import (
	"context"
	"fmt"
	"log"
	"sort"
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
	// addUps := make([]model.Score, 0)
	arrgedAddUps := make([][]model.Score, 0)
	ch := make(chan model.JudgeScore)             // 創建channel
	ch_addUps := make(chan *[]model.Score)        // addUps channel
	ch_addUpResponse := make(chan *[]model.Score) // addUps channel
	go addUpScores(ch_addUps, ch_addUpResponse, (*judges)[0].RowNum, len(*judges))
	for i := 0; i < len(*judges); i++ {
		go getAllScores(&((*judges)[i]), ch, ch_addUps) //同步取得多位judge的所有分數，並使用channel回傳。
	}
	result := make([]model.JudgeScore, 0) // 創造空slice並依序把channel收到的回傳值插入
	for x := 0; x < len(*judges); x++ {   // 有幾位judge就接收channel幾次
		result = append(result, <-ch)
	}
	addUps := <-ch_addUpResponse
	sortByPoint(addUps)
	rankScores(addUps)
	sortByRowandNum(addUps)
	arrangeAddUps(addUps, &arrgedAddUps, (*judges)[0].RowNum)
	c.JSON(200, gin.H{ //回傳
		"message":   true,
		"allScores": result,
		"addUps":    arrgedAddUps,
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
func getAllScores(judge *model.Judge, ch_score chan model.JudgeScore, ch_addUps chan *[]model.Score) {
	judgeID := judge.ID.Hex()
	collection := database.Db.Collection("scores")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	findOptions := options.Find().SetSort(bson.D{{"point", -1}})
	cursor, err := collection.Find(ctx, bson.M{"judgeID": judgeID}, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	var scores []model.Score
	if err = cursor.All(ctx, &scores); err != nil {
		log.Fatal(err)
	}
	fmt.Println(judgeID)
	rankScores(&scores)
	sortByRowandNum(&scores)
	// addUpScores(addUps, &scores, judge.RowNum)
	ch_addUps <- &scores
	// 如果沒有分數就回傳空陣列
	if len(scores) > 0 {
		var arrangedScores = make([][]model.Score, 0)
		reArrangeScores(scores, &arrangedScores, judge.RowNum)
		var judgeScore = model.JudgeScore{
			ID:     judge.ID,
			Name:   judge.Name,
			RowNum: judge.RowNum,
			Scores: arrangedScores,
		}
		ch_score <- judgeScore
	} else {
		var judgeScore = model.JudgeScore{
			ID:     judge.ID,
			Name:   judge.Name,
			RowNum: judge.RowNum,
			Scores: [][]model.Score{},
		}
		ch_score <- judgeScore
	}
}

// 按照 row、 number 重新排序所有分數，中間有空的分數直接補入空值
func reArrangeScores(scores []model.Score, arrangedScores *[][]model.Score, rowNum int) {
	var scoreArray *[][]model.Score = arrangedScores
	var subArray []model.Score
	r, n := 1, 1

	for i := 0; i < len(scores); {
		if scores[i].Number == n && scores[i].Row == r { //該號碼該排有分數就插入
			currentScore := model.Score{
				ID:      scores[i].ID,
				Row:     scores[i].Row,
				Number:  scores[i].Number,
				Point:   scores[i].Point,
				Rank:    scores[i].Rank,
				IsEmpty: false,
			}
			subArray = append(subArray, currentScore)
			// fmt.Println("append 1, currentRow: ", r, " currentNum: ", n)
			r++
			i++
		} else { //該號碼該排沒分數就插入nil
			subArray = append(subArray, model.Score{IsEmpty: true})
			// fmt.Println("append empty, currentRow: ", r, " currentNum: ", n)
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
			// fmt.Println("last")
			if r != 1 && r <= rowNum {
				for r <= rowNum {
					subArray = append(subArray, model.Score{IsEmpty: true})
					r++
				}
				*scoreArray = append(*scoreArray, subArray)
			}
		}
	}
}

// 排序分組addUps
func arrangeAddUps(addUps *[]model.Score, arragedAddUps *[][]model.Score, rowNum int) {
	subArray := make([]model.Score, 0)
	counter := 1
	for i, score := range *addUps {
		if counter > rowNum {
			*arragedAddUps = append(*arragedAddUps, subArray)
			subArray = nil
			counter = 1
		}
		subArray = append(subArray, score)
		counter += 1
		if i == len(*addUps)-1 && len(subArray) > 0 {
			for j := 0; j < rowNum-len(subArray); {
				subArray = append(subArray, model.Score{IsEmpty: true})
			}
			*arragedAddUps = append(*arragedAddUps, subArray)
		}
	}
}

func rankScores(scores *[]model.Score) {
	for i, _ := range *scores {
		if i > 0 && (*scores)[i].Point == (*scores)[i-1].Point {
			(*scores)[i].Rank = (*scores)[i-1].Rank
		} else {
			(*scores)[i].Rank = i + 1
		}
	}

}

func sortByRowandNum(scores *[]model.Score) {
	sort.Slice(*scores, func(i, j int) bool {
		if (*scores)[i].Number == (*scores)[j].Number {
			return (*scores)[i].Row < (*scores)[j].Row
		} else {
			return (*scores)[i].Number < (*scores)[j].Number
		}
	})
}

func sortByPoint(scores *[]model.Score) {
	sort.Slice(*scores, func(i, j int) bool {
		return (*scores)[i].Point > (*scores)[j].Point
	})
}

func addUpScores(ch chan *[]model.Score, ch_response chan *[]model.Score, rowNum int, judgeNum int) {
	addUps := make([]model.Score, 0)
	allScores := make([][]model.Score, 0)
	for i := 0; i < judgeNum; { // 等待接收所有judges 的scores
		// 分數進來後先插入必要的空白分數

		scores := <-ch

		r, n := 1, 1
		arrangedScores := make([]model.Score, 0) // 加入空白分數後的分數

		for i := 0; i < len(*scores); {
			if r > rowNum {
				r = 1
				n += 1
			}
			if (*scores)[i].Row == r && (*scores)[i].Number == n {
				arrangedScores = append(arrangedScores, (*scores)[i])
				i += 1
			} else {
				emptyScore := model.Score{
					Row:     r,
					Number:  n,
					IsEmpty: true,
				}
				arrangedScores = append(arrangedScores, emptyScore)
			}
			r += 1
		}
		allScores = append(allScores, arrangedScores)
		i++
	}

	for _, scores := range allScores { // 依次處理所有的scores 一個個加進addUps
		if len(addUps) > 0 {
			for i, score := range scores {
				if i > len(addUps)-1 { // 若新來的分數長度超過addUps就直接插入
					addUps = append(addUps, model.Score{Row: score.Row, Number: score.Number, Point: score.Point, IsEmpty: score.IsEmpty})
				} else {
					// fmt.Print("\nscore ", score.Number, score.Row, "   ")
					// fmt.Println("addups ", (addUps)[i].Number, (addUps)[i].Row)

					(addUps)[i].Point += score.Point
					(addUps)[i].Point = (addUps)[i].Point * 100 / 100
					(addUps)[i].IsEmpty = (addUps)[i].IsEmpty || score.IsEmpty // 判定是否為empty
				}
			}
		} else {
			for _, score := range scores {
				// fmt.Print("appending first set ", score.Number, score.Row)
				addUps = append(addUps, model.Score{Row: score.Row, Number: score.Number, Point: score.Point, IsEmpty: score.IsEmpty})
			}
		}
	}
	ch_response <- &addUps
}
