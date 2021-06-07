package model

import "go.mongodb.org/mongo-driver/bson/primitive"

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
