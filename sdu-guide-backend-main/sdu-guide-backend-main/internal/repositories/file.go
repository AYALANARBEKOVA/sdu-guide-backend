package repositories

import (
	"context"
	"sdu-guide/internal/ai"
	"sdu-guide/internal/structures"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type File struct {
	db         *mongo.Database
	ai         *ai.AI
	collection *mongo.Collection
}

func newFileRepo(db *mongo.Database, ai *ai.AI, collection *mongo.Collection) *File {
	return &File{db: db, ai: ai, collection: collection}
}

func (m *File) CreateFile(file structures.File) error {

	file.ID = m.ai.Next("files")

	if _, err := m.collection.InsertOne(context.Background(), file); err != nil {
		return err
	}
	return nil
}

func (m *File) UpadteXLSX(file structures.File) error {

	if _, err := m.collection.UpdateOne(context.Background(), bson.M{"_id": file.ID}, bson.M{"$set": file}); err != nil {
		return err
	}
	return nil
}

func (m *File) GetFile(hash string) (structures.File, error) {
	result := structures.File{}
	if err := m.collection.FindOne(context.Background(), bson.M{"hash": hash}).Decode(&result); err != nil {
		return structures.File{}, err
	}
	return result, nil
}

func (m *File) Delete(hash string) error {

	if _, err := m.collection.DeleteOne(context.Background(), bson.M{"hash": hash}); err != nil {
		return err
	}
	return nil
}
