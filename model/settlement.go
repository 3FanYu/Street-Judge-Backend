package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Settlement struct {
	JudgeScore []JudgeScore `json:"judgeScore,omitempty" bson:"judgeScore,omitempty"`
}
type JudgeScore struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	EventID string             `json:"eventID,omitempty" bson:"eventID,omitempty"`
	Name    string             `json:"name,omitempty" bson:"name,omitempty"`
	RowNum  int                `json:"rowNum,omitempty" bson:"rowNum,omitempty"`
	Scores  []Score            `json:"scores" bson:"scores,omitempty"`
}
