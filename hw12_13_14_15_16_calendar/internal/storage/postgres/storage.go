package postgres

import (
	"context"
	"time"

	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/internal/app"
)

type Storage struct { // TODO
}

func New() *Storage {
	return &Storage{}
}

var _ app.Storage = &Storage{} // check the interface

func (s *Storage) Connect(ctx context.Context) error {
	// TODO
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	// TODO
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

func (s *Storage) ListEventsForMonth(month time.Month) []app.Event {
	//TODO implement me
	panic("implement me")
}
