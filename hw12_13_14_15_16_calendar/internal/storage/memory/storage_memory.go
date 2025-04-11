package storagememory

import (
	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/pkg/logger"
	"sync"
	"time"

	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/internal/app"
)

type Storage struct {
	mu     sync.RWMutex
	logg   *logger.Logger
	events map[string]app.Event
}

var _ app.Storage = &Storage{}

func New(logg *logger.Logger) *Storage {
	logg.Info("initial storage in-memory")
	return &Storage{
		logg:   logg,
		events: make(map[string]app.Event),
	}
}

func (s *Storage) Connect() error {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) Close() error {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) CreateEvent(event app.Event) (app.Event, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) UpdateEvent(event app.Event) (app.Event, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) DeleteEvent(event app.Event) error {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) ListEventsForDay(day time.Time) []app.Event {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) ListEventsForWeek(week time.Time) []app.Event {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) ListEventsForMonth(month time.Time) []app.Event {
	//TODO implement me
	panic("implement me")
}
