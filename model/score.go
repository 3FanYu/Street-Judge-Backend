package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Score struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	EventID string             `json:"eventID,omitempty" bson:"eventID,omitempty"`
	Row     int                `json:"row,omitempty" bson:"row,omitempty"`
	Name    string             `json:"name,omitempty" bson:"name,omitempty"`
	Point   float32            `json:"point,omitempty" bson:"point,omitempty"`
}
