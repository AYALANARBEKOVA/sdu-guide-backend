package services

import (
	"mime/multipart"
	"net/http"
	"sdu-guide/internal/repositories"
	"sdu-guide/internal/structures"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Service struct {
	UserService
	FileServices
	RoomService
	EventService
}

func NewService(repo *repositories.Repository) *Service {
	return &Service{
		UserService:  newUserService(repo),
		FileServices: newFileService(repo),
		RoomService:  newRoomService(repo),
		EventService: newEventService(repo),
	}
}

type UserService interface {
	Get(id int64) (structures.User, error)
	GetBy(field string, value interface{}) (structures.User, error)
	Create(userRegister structures.UserRegister) error
	Upadte(user structures.User) error
	Logout(cookie *http.Cookie) error
	LoginToSystem(login structures.Login) (*http.Cookie, structures.Session, error)
	GetUserFromSession(session structures.Session) (structures.User, error)
}

type FileServices interface {
	StoreFile(file multipart.File, header *multipart.FileHeader) (string, error)
	GetFilebyHash(hash string) (structures.File, error)
}

type RoomService interface {
	CreateRoom(room structures.Room) error
	GetRoom(id int64) (structures.Room, error)
	UpdateRoom(room structures.Room) error
	GetAll(filter bson.M) ([]structures.Room, error)
	Delete(id uint64) error
	GetScheduleBySEF(sef string) (structures.File, error)
}

type EventService interface {
	CreateEvent(event structures.Event) error
	GetEvent(id int64) (structures.Event, error)
	UpdateEvent(event structures.Event) error
	GetAll(filter structures.Filter) ([]structures.Event, error)
	Delete(id uint64) error
}
