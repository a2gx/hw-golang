package app

import (
	"context"
	"time"
)

type Application interface {
	CreateEvent(event Event) (Event, error)
	UpdateEvent(event Event) (Event, error)
	DeleteEvent(event Event) error

	ListEventsForDay(day time.Time) []Event
	ListEventsForWeek(week time.Time) []Event
	ListEventsForMonth(month time.Time) []Event
}

type Storage interface {
	Application

	Connect() error
	Close() error
}

type Server interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}
