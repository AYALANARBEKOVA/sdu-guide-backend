package repositories

import (
	"context"
	"sdu-guide/internal/ai"
	"sdu-guide/internal/structures"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Event struct {
	db         *mongo.Database
	ai         *ai.AI
	collection *mongo.Collection
}

func newEventRepo(db *mongo.Database, ai *ai.AI, collection *mongo.Collection) *Event {
	repo := &Event{db: db, ai: ai, collection: collection}
	return repo
}

func (m *Event) Create(event structures.Event) error {

	event.ID = int64(m.ai.Next("events"))
	if _, err := m.collection.InsertOne(context.Background(), event); err != nil {
		return err
	}
	return nil
}

func (m *Event) Update(event structures.Event) error {

	if _, err := m.collection.UpdateByID(context.Background(), event.ID, bson.M{"$set": event}); err != nil {
		return err
	}
	return nil
}

func (m *Event) Delete(id int64) error {
	if _, err := m.collection.DeleteOne(context.Background(), bson.M{"_id": id}); err != nil {
		return err
	}

	return nil
}

func (m *Event) Get(id int64) (structures.Event, error) {
	result := structures.Event{}

	if err := m.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&result); err != nil {
		return structures.Event{}, err
	}

	return result, nil
}

func (m *Event) GetBy(obj bson.M) (structures.Event, error) {
	result := structures.Event{}

	if err := m.collection.FindOne(context.Background(), obj).Decode(&result); err != nil {
		return structures.Event{}, err
	}

	return result, nil
}

func (m *Event) GetAll(filter structures.Filter) ([]structures.Event, error) {
	result := []structures.Event{}
	opt := options.Find().SetSort(bson.M{"_id": -1})
	if filter.Limit > 0 {
		opt.SetLimit(filter.Limit)
	}
	cur, err := m.collection.Find(context.Background(), filter.Request, opt)
	if err != nil {
		return result, err
	}

	if err := cur.All(context.Background(), &result); err != nil {
		return result, err
	}
	return result, nil
}

func (m *Event) MarkPastEventsAsEnded() error {
	now := time.Now()

	filter := bson.M{
		"date":  bson.M{"$lt": now}, // События, где дата меньше текущей
		"ended": false,              // Только те, которые еще не помечены как завершенные
	}

	update := bson.M{
		"$set": bson.M{"ended": true},
	}

	_, err := m.collection.UpdateMany(context.Background(), filter, update)
	return err
}
