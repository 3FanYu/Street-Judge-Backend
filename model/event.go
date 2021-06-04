package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Event struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name,omitempty" bson:"name,omitempty"`
	Owner    string             `json:"owner,omitempty" bson:"owner,omitempty"`
	Password string             `json:"password,omitempty" bson:"password,omitempty"`
	Judges   []string           `json:"judges,omitempty" bson:"judges,omitempty"`
	RowNum   int                `json:"rowNum,omitempty" bson:"rowNum,omitempty"`
}
