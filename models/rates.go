package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Rate struct {
	MarketName     string    `json:"MarketName" bson:"market-name"`
	High           float32   `json:"High" bson:"high"`
	Low            float32   `json:"Low" bson:"low"`
	Volume         float32   `json:"Volume" bson:"volume"`
	Last           float32   `json:"Last" bson:"last"`
	BaseVolume     float32   `json:"BaseVolume" bson:"base-volume"`
	Bid            float32   `json:"Bid" bson:"bid"`
	Ask            float32   `json:"Ask" bson:"ask"`
	PrevDay        float32   `json:"PrevDay" bson:"prev-day"`
	OpenBuyOrders  uint      `json:"OpenBuyOrders" bson:"open-buy-orders"`
	OpenSellOrders uint      `json:"OpenBuyOrders" bson:"open-sell-orders"`
	TimeStamp      time.Time `json:"TimeStamp" bson:"time-stamp"`
	Created        time.Time `json:"Created" bson:"created"`
}

type Rates struct {
	Id        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Rate      []Rate             `json:"Rates" bson:"rate"`
	TimeStamp time.Time          `json:"TimeStamp" bson:"time-stamp"`
}
