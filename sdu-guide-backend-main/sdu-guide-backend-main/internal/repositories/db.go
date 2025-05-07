package repositories

import (
	"sdu-guide/internal/logger"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

func NewDB() (*mongo.Database, error) {

	client, err := mongo.Connect()

	if err != nil {
		logger.Error.Println(err)
		return nil, err
	}

	db := client.Database("sud-guide-test")
	return db, nil
}
