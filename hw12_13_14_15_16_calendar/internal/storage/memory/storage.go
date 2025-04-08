package memory

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/pkg/logger"

	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/internal/app"
)

type Storage struct {
	mu     sync.RWMutex
	events map[string]app.Event
}

var _ app.Storage = &Storage{} // check the interface

func New(logg *logger.Logger) *Storage {
	logg.Info("start storage -> memory")
	return &Storage{
		events: make(map[string]app.Event),
	}
}

// CreateEvent implements app.Storage.
func (s *Storage) CreateEvent(event app.Event) (app.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.events[event.ID]; exists {
		return app.Event{}, errors.New("event already exists")
	}

	s.events[event.ID] = event
	return event, nil
}

// UpdateEvent implements app.Storage.
func (s *Storage) UpdateEvent(event app.Event) (app.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.events[event.ID]; !exists {
		return app.Event{}, errors.New("event not found")
	}

	s.events[event.ID] = event
	return event, nil
}

// DeleteEvent implements app.Storage.
func (s *Storage) DeleteEvent(event app.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.events[event.ID]; !exists {
		return errors.New("event not found")
	}

	delete(s.events, event.ID)
	return nil
}

// ListEventsForDay implements app.Storage.
func (s *Storage) ListEventsForDay(date time.Time) []app.Event {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var events []app.Event
	for _, event := range s.events {
		if event.Datetime.Year() == date.Year() &&
			event.Datetime.Month() == date.Month() &&
			event.Datetime.Day() == date.Day() {
			events = append(events, event)
		}
	}
	return events
}

// ListEventsForWeek implements app.Storage.
func (s *Storage) ListEventsForWeek(date time.Time) []app.Event {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var events []app.Event
	weekStart := date.AddDate(0, 0, -int(date.Weekday()))
	weekEnd := weekStart.AddDate(0, 0, 7)

	for _, event := range s.events {
		if event.Datetime.After(weekStart) && event.Datetime.Before(weekEnd) {
			events = append(events, event)
		}
	}
	return events
}

// ListEventsForMonth implements app.Storage.
func (s *Storage) ListEventsForMonth(month time.Month) []app.Event {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var events []app.Event
	for _, event := range s.events {
		if event.Datetime.Month() == month {
			events = append(events, event)
		}
	}
	return events
}

func (s *Storage) Close(_ context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.events = make(map[string]app.Event)
	return nil
}

func (s *Storage) Connect(_ context.Context) error {
	return nil
}
