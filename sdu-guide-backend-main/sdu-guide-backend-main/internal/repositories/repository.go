package repositories

import (
	"sdu-guide/internal/ai"
	"sdu-guide/internal/structures"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository struct {
	UserRepo
	RoomRepo
	FileRepo
	EventRepo
}

func NewRepository(db *mongo.Database, ai *ai.AI) *Repository {
	return &Repository{
		UserRepo:  newUserRepo(db, ai, db.Collection("users")),
		RoomRepo:  newRoomRepo(db, ai, db.Collection("rooms")),
		FileRepo:  newFileRepo(db, ai, db.Collection("files")),
		EventRepo: newEventRepo(db, ai, db.Collection("events")),
	}
}

type UserRepo interface {
	Get(id int64) (structures.User, error)
	GetBy(field string, value interface{}) (structures.User, error)
	Create(user structures.User) error
	Upadte(user structures.User) error
}

type RoomRepo interface {
	Create(room structures.Room) error
	GetBy(field string, value interface{}) (structures.Room, error)
	GetAll(filter bson.M) ([]structures.Room, error)
	Update(room structures.Room) error
	Delete(id uint64) error
}

type FileRepo interface {
	CreateFile(file structures.File) error
	UpadteXLSX(file structures.File) error
	GetFile(hash string) (structures.File, error)
	Delete(hash string) error
}

type EventRepo interface {
	Create(event structures.Event) error
	Update(event structures.Event) error
	Delete(id int64) error
	Get(id int64) (structures.Event, error)
	GetAll(filter structures.Filter) ([]structures.Event, error)
	MarkPastEventsAsEnded() error
	GetBy(obj bson.M) (structures.Event, error)
}
