package model

import (
	"context"
	"time"

	"github.com/3FanYu/Judges321-backend/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Event struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name,omitempty" bson:"name,omitempty"`
	Owner    string             `json:"owner,omitempty" bson:"owner,omitempty"`
	Password string             `json:"password,omitempty" bson:"password,omitempty"`
	Judges   []Judge            `json:"judges,omitempty" bson:"judges,omitempty"`
}

func (event *Event) CreateEvent() (*mongo.InsertOneResult, error) {
	collection := database.Db.Collection("events")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	res, err := collection.InsertOne(ctx, &event)

	return res, err
}

func (event *Event) GetEvent(ID *primitive.ObjectID) {
	collection := database.Db.Collection("events")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	res := collection.FindOne(ctx, bson.M{"_id": *ID}, options.FindOne().SetProjection(bson.M{"password": 0}))
	// var result Event
	res.Decode(event)
}
