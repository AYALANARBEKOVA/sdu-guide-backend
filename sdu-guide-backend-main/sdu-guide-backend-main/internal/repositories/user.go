package repositories

import (
	"context"
	"sdu-guide/internal/ai"
	"sdu-guide/internal/structures"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type User struct {
	db         *mongo.Database
	ai         *ai.AI
	collection *mongo.Collection
}

func newUserRepo(db *mongo.Database, ai *ai.AI, collection *mongo.Collection) *User {
	return &User{db: db, ai: ai, collection: collection}
}

func (m *User) Get(id int64) (structures.User, error) {
	var user structures.User
	if err := m.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user); err != nil {
		return structures.User{}, err
	}
	return user, nil
}

func (m *User) GetBy(field string, value interface{}) (structures.User, error) {
	var user structures.User
	if err := m.collection.FindOne(context.Background(), bson.M{field: value}).Decode(&user); err != nil {
		return structures.User{}, err
	}
	return user, nil
}

func (m *User) Create(user structures.User) error {

	user.ID = m.ai.Next("users")

	if _, err := m.collection.InsertOne(context.Background(), user); err != nil {
		return err
	}
	return nil
}

func (m *User) Upadte(user structures.User) error {

	if _, err := m.collection.UpdateOne(context.Background(), bson.M{"_id": user.ID}, bson.M{"$set": user}); err != nil {
		return err
	}
	return nil
}
