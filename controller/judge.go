package controller

import (
	"context"
	"log"
	"time"

	"github.com/3FanYu/Judges321-backend/database"
	"github.com/3FanYu/Judges321-backend/model"
	"github.com/gin-gonic/gin"
)

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
