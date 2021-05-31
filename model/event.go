package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Event struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name,omitempty" bson:"name,omitempty"`
	Owner     string             `json:"owner,omitempty" bson:"owner,omitempty"`
	Password  string             `json:"password,omitempty" bson:"password,omitempty"`
	JudgeName []string           `json:"judgeName,omitempty" bson:"judgeName,omitempty"`
	Judges    []Judge            `json:"judges,omitempty" bson:"judges,omitempty"`
	RowNum    int                `json:"rowNum,omitempty" bson:"rowNum,omitempty"`
}

type Judge struct {
	Name   string  `json:"name,omitempty" bson:"name,omitempty"`
	Scores []Score `json:"scores,omitempty" bson:"scores,omitempty"`
}

type Score struct {
	Point float32 `json:"point,omitempty" bson:"point,omitempty"`
}

func NewScore() *Score {
	return &Score{
		Point: 0,
	}
}
