package structures

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Event struct {
	ID        int64     `json:"id" bson:"_id"`
	Name      string    `json:"name" bson:"name"`
	Date      time.Time `json:"date" bson:"date"`
	Place     string    `json:"place" bson:"place"`
	StartTime string    `json:"startTime" bson:"startTime"`
	EndTime   string    `json:"endTime" bson:"endTime"`
	Ended     bool      `json:"ended" bson:"ended"`
	Hash      string    `json:"hash" bson:"hash"`
	ShortName string    `json:"shortName" bson:"shortName"`
}

type Filter struct {
	Limit   int64
	Request bson.M
}
