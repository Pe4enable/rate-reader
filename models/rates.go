package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Rate struct {
	MarketName     string  `json:"MarketName"`
	High           float64 `json:"High"`
	Low            float64 `json:"Low"`
	Volume         float64 `json:"Volume"`
	Last           float64 `json:"Last"`
	BaseVolume     float64 `json:"BaseVolume"`
	Bid            float64 `json:"Bid"`
	Ask            float64 `json:"Ask"`
	OpenBuyOrders  uint     `json:"OpenBuyOrders"`
	OpenSellOrders uint     `json:"OpenSellOrders"`
	PrevDay        float64 `json:"PrevDay"`
}

type Rates struct {
	Id        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Rates     []Rate             `json:"Rates" bson:"rate"`
	TimeStamp time.Time          `json:"TimeStamp" bson:"time-stamp"`
}