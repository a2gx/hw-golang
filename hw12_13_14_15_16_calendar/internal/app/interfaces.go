package app

import (
	"context"
	"time"
)

type Storage interface {
	Connect() error
	Close() error

	CreateEvent(event Event) (Event, error)
	UpdateEvent(event Event) (Event, error)
	DeleteEvent(event Event) error

	ListEventsForDay(day time.Time) []Event
	ListEventsForWeek(week time.Time) []Event
	ListEventsForMonth(month time.Time) []Event
}

type Server interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}
