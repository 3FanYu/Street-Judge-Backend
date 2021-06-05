package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Score struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	JudgeID string             `json:"judgeID,omitempty" bson:"judgeID,omitempty"`
	Row     int                `json:"row,omitempty" bson:"row,omitempty"`
	Number  int                `json:"number,omitempty" bson:"number,omitempty"`
	Point   float32            `json:"point,omitempty" bson:"point,omitempty"`
}
