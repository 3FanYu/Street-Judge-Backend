package model

import (
	"context"
	"time"

	"github.com/3FanYu/Judges321-backend/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Judge struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	EventID string             `json:"eventID,omitempty" bson:"eventID,omitempty"`
	Name    string             `json:"name,omitempty" bson:"name,omitempty"`
	Names   []string           `json:"names,omitempty" bson:"names,omitempty"`
	RowNum  int                `json:"rowNum,omitempty" bson:"rowNum,omitempty"`
}

// type Judges struct {
// 	EventID string   `json:"eventID,omitempty" bson:"eventID,omitempty"`
// 	RowNum  int      `json:"rowNum,omitempty" bson:"rowNum,omitempty"`
// }

type JudgeInfo struct {
	EventName  string `json:"eventName,omitempty" bson:"eventName,omitempty"`
	EventOwner string `json:"eventOwner,omitempty" bson:"eventOwner,omitempty"`
	JudgeName  string `json:"judgeName,omitempty" bson:"judgeName,omitempty"`
	RowNum     int    `json:"rowNum,omitempty" bson:"rowNum,omitempty"`
}

func (judges *Judge) CreateJudge() (*mongo.InsertManyResult, error) {
	data := judges.makeInterface()
	collection := database.Db.Collection("judges")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	res, err := collection.InsertMany(ctx, *data)
	if err != nil {
		res = nil
	}
	return res, err
}

func (judges *Judge) makeInterface() *[]interface{} {
	//資料放進[]interfce{}
	var tmp []interface{}
	for _, v := range judges.Names {
		var judge = Judge{
			EventID: judges.EventID,
			Name:    v,
			RowNum:  judges.RowNum,
		}
		tmp = append(tmp, judge)
	}
	return &tmp
}

func (judge *Judge) GetJudge(ID *primitive.ObjectID) {
	collection := database.Db.Collection("judges")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	res := collection.FindOne(ctx, bson.M{"_id": ID})
	res.Decode(judge)
}
