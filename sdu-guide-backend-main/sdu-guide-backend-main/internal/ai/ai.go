package ai

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type AI struct {
	c *mongo.Collection
}

func NewAI(db *mongo.Database) AI {
	return AI{c: db.Collection("ai")}
}

func (a *AI) Next(name string) uint64 {
	seqFieldName := "seq"
	idFieldName := "id"
	result := bson.M{}
	if err := a.c.FindOneAndUpdate(context.Background(),
		bson.M{idFieldName: name},
		bson.M{"$set": bson.M{idFieldName: name},
			"$inc": bson.M{seqFieldName: 1}},
		options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After),
	).Decode(&result); err != nil {
		fmt.Println("Autoincrement error(1):", err.Error())
	}

	sec, ok := result[seqFieldName].(int)
	if ok {
		return uint64(sec)
	}

	sec2, ok := result[seqFieldName].(int64)
	if ok {
		return uint64(sec2)
	}

	sec3, _ := result[seqFieldName].(int32)
	return uint64(sec3)
}
