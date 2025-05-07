package services

import (
	"errors"
	"sdu-guide/internal/repositories"
	"sdu-guide/internal/structures"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Room struct {
	repo *repositories.Repository
}

func newRoomService(repo *repositories.Repository) *Room {
	return &Room{repo: repo}
}

func (m *Room) CreateRoom(room structures.Room) error {
	if len(room.Block) == 0 {
		return errors.New("choose room block")
	}
	if room.Number <= 0 {
		return errors.New("choose room number")
	}
	room.Updated = time.Now()
	if len(strings.TrimSpace(room.SEF)) == 0 {
		room.SEF = strings.ToLower(room.Block + "-" + strconv.Itoa(int(room.Number)))
	}
	return m.repo.RoomRepo.Create(room)
}

func (m *Room) GetRoom(id int64) (structures.Room, error) {
	return m.repo.RoomRepo.GetBy("_id", id)
}

func (m *Room) UpdateRoom(room structures.Room) error {

	_, err := m.repo.RoomRepo.GetBy("_id", room.ID)
	if err != nil {
		return err
	}

	if len(strings.TrimSpace(room.SEF)) == 0 {
		room.SEF = strings.ToLower(room.Block + "-" + strconv.Itoa(int(room.Number)))
	}
	room.Updated = time.Now()

	if err := m.repo.RoomRepo.Update(room); err != nil {
		return err
	}

	return nil
}

func (m *Room) GetAll(filter bson.M) ([]structures.Room, error) {

	return m.repo.RoomRepo.GetAll(filter)

}

func (m *Room) Delete(id uint64) error {
	return m.repo.RoomRepo.Delete(id)
}

func (m *Room) GetScheduleBySEF(sef string) (structures.File, error) {

	room, err := m.repo.RoomRepo.GetBy("sef", sef)
	if err != nil {
		return structures.File{}, errors.New("room don't exist")
	}
	if room.ScheduleHash == "" {
		return structures.File{}, errors.New("room don't have schedule")
	}
	return m.repo.FileRepo.GetFile(room.ScheduleHash)
}
