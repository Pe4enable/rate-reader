package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Rate struct {
	Symbol       string  `json:"Symbol"`
	High         float64 `json:"High,string"`
	HighChange   float64 `json:"HighChange"`
	Low          float64 `json:"Low,string"`
	LowChange    float64 `json:"LowChange"`
	Volume       float64 `json:"Volume,string"`
	VolumeChange float64 `json:"VolumeChange"`
}

type Rates struct {
	Id        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Rates     []Rate             `json:"Rates" bson:"rate"`
	TimeStamp time.Time          `json:"TimeStamp" bson:"time-stamp"`
}
