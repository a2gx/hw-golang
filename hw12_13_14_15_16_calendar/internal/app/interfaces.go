package app

import (
	"context"
	"time"
)

type Application interface {
	Create(event Event) (Event, error)
	Update(event Event) (Event, error)
	Delete(event Event) error

	ListEventsForDay(day time.Time) []Event
	ListEventsForWeek(week time.Time) []Event
	ListEventsForMonth(month time.Time) []Event
}

type Storage interface {
	Connect() error
	Close() error

	CreateEvent(event Event) (Event, error)
	UpdateEvent(event Event) (Event, error)
	DeleteEvent(event Event) error

	ListEventsForDay(day time.Time) []Event
	ListEventsForWeek(week time.Time) []Event
	ListEventsForMonth(month time.Time) []Event

	GetById(eventId string) (Event, error)
}

type Server interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}
