package repositories

import (
	"context"
	"sdu-guide/internal/ai"
	"sdu-guide/internal/logger"
	"sdu-guide/internal/structures"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Room struct {
	db         *mongo.Database
	ai         *ai.AI
	collection *mongo.Collection
}

func newRoomRepo(db *mongo.Database, ai *ai.AI, collection *mongo.Collection) *Room {
	repo := &Room{db: db, ai: ai, collection: collection}
	if err := repo.ensureIndexes(); err != nil {
		logger.Error.Printf("Ошибка при создании индексов: %v", err)
	}
	return repo
}

func (m *Room) ensureIndexes() error {
	ctx := context.Background()
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "sef", Value: 1}}, // Индекс по sef
			Options: options.Index().SetUnique(true),
		},
	}

	_, err := m.collection.Indexes().CreateMany(ctx, indexes)
	return err
}

func (m *Room) Create(room structures.Room) error {
	room.ID = m.ai.Next("rooms")
	if _, err := m.collection.InsertOne(context.Background(), room); err != nil {
		return err
	}
	return nil
}

func (m *Room) GetBy(field string, value interface{}) (structures.Room, error) {
	room := structures.Room{}
	if err := m.collection.FindOne(context.Background(), bson.M{field: value}).Decode(&room); err != nil {
		return structures.Room{}, err
	}
	return room, nil
}

func (m *Room) GetAll(filter bson.M) ([]structures.Room, error) {

	cur, err := m.collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	result := []structures.Room{}
	if err := cur.All(context.Background(), &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (m *Room) Update(room structures.Room) error {

	if _, err := m.collection.UpdateOne(context.Background(), bson.M{"_id": room.ID}, bson.M{"$set": room}); err != nil {
		return err
	}
	return nil
}

func (m *Room) Delete(id uint64) error {

	if _, err := m.collection.DeleteOne(context.Background(), bson.M{"_id": id}); err != nil {
		return err
	}
	return nil
}
