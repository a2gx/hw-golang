package app

import (
	"context"
	"time"
)

type Storage interface {
	Connect(ctx context.Context) error
	Close(ctx context.Context) error

	CreateEvent(event Event) (Event, error)
	UpdateEvent(event Event) (Event, error)
	DeleteEvent(event Event) error

	ListEventsForDay(day time.Time) ([]Event, error)
	ListEventsForWeek(week time.Time) ([]Event, error)
	ListEventsForMonth(month time.Month) ([]Event, error)
}
