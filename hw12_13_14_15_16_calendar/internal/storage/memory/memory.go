package memory

import (
	"context"
	"time"

	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/internal/app"
)

type Storage struct{}

var _ app.Storage = &Storage{}

func New() *Storage {
	return &Storage{}
}

func (m *Storage) Connect(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (m *Storage) Close(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (m *Storage) CreateEvent(event app.Event) (app.Event, error) {
	//TODO implement me
	panic("implement me")
}

func (m *Storage) UpdateEvent(event app.Event) (app.Event, error) {
	//TODO implement me
	panic("implement me")
}

func (m *Storage) DeleteEvent(event app.Event) error {
	//TODO implement me
	panic("implement me")
}

func (m *Storage) ListEventsForDay(day time.Time) ([]app.Event, error) {
	//TODO implement me
	panic("implement me")
}

func (m *Storage) ListEventsForWeek(week time.Time) ([]app.Event, error) {
	//TODO implement me
	panic("implement me")
}

func (m *Storage) ListEventsForMonth(month time.Month) ([]app.Event, error) {
	//TODO implement me
	panic("implement me")
}
