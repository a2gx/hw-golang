package app

import (
	"context"
	"time"
)

type Application interface {
	Create(event Event) (Event, error)
	Update(event Event) (Event, error)
	Delete(event Event) error

	EventsInInterval(date time.Time, days int) []Event
}

type Storage interface {
	Connect() error
	Close() error

	CreateEvent(event Event) (Event, error)
	UpdateEvent(event Event) (Event, error)
	DeleteEvent(event Event) error

	GetByID(eventID string) (Event, error)
	FilterByInterval(st, fn time.Time) []Event

	FetchEventsToNotify() ([]Event, error)
	DeleteOldEvents() error
}

type Server interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}
