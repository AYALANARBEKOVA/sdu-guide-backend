package services

import (
	"sdu-guide/internal/logger"
	"sdu-guide/internal/repositories"
	"sdu-guide/internal/structures"

	"github.com/robfig/cron/v3"
)

type Event struct {
	repo *repositories.Repository
}

func newEventService(repo *repositories.Repository) *Event {
	service := &Event{repo: repo}

	service.setupCron()
	return service
}

func (m *Event) CreateEvent(event structures.Event) error {

	return m.repo.EventRepo.Create(event)
}

func (m *Event) GetEvent(id int64) (structures.Event, error) {
	return m.repo.EventRepo.Get(id)
}

func (m *Event) UpdateEvent(event structures.Event) error {

	_, err := m.repo.EventRepo.Get(event.ID)
	if err != nil {
		return err
	}

	if err := m.repo.EventRepo.Update(event); err != nil {
		return err
	}

	return nil
}

func (m *Event) GetAll(filter structures.Filter) ([]structures.Event, error) {

	return m.repo.EventRepo.GetAll(filter)

}

func (m *Event) Delete(id uint64) error {
	return m.repo.EventRepo.Delete(int64(id))
}

func (m *Event) setupCron() {
	c := cron.New()

	_, err := c.AddFunc("0 0 * * * ", func() { // Функция будет срабатывает каждый день в  00:00
		err := m.repo.EventRepo.MarkPastEventsAsEnded()
		if err != nil {
			logger.Error.Println("Can't update data:", err)
		} else {
			logger.Info.Println("Old data was set as ended in database")
		}

	})
	if err != nil {
		logger.Error.Printf("Cron error: %v\n", err)
		return
	}

	c.Start()
}
