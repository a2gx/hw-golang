package storage_sql

import (
	"time"

	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/internal/app"
	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/pkg/logger"
)

type Storage struct {
	logg *logger.Logger
}

var _ app.Storage = &Storage{}

func New(logg *logger.Logger) *Storage {
	return &Storage{
		logg: logg,
	}
}

func (s *Storage) Connect() error {
	s.logg.Debug("storage SQL connected")
	return nil
}

func (s *Storage) Close() error {
	s.logg.Debug("storage SQL closed")
	return nil
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
