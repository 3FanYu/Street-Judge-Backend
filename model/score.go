package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Score struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	JudgeID string             `json:"judgeID,omitempty" bson:"judgeID,omitempty"`
	Row     int                `json:"row" bson:"row" binding:"required"`
	Number  int                `json:"number" bson:"number" binding:"required"`
	Point   float32            `json:"point" bson:"point" `
	IsEmpty bool               `json:"isEmpty" bson:"isEmpty,omitempty" `
	Rank    int                `json:"rank,omitempty" bson:"rank,omitempty" `
}
