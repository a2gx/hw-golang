package memory

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/internal/app"
	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/pkg/logger"
)

type Storage struct {
	mu     sync.RWMutex
	events map[string]app.Event
}

var _ app.Storage = &Storage{}

func New(logg *logger.Logger) *Storage {
	logg.Info("start storage -> memory")
	return &Storage{
		events: make(map[string]app.Event),
	}
}

// Close implements app.Storage.
func (s *Storage) Close(_ context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.events = make(map[string]app.Event)

	return nil
}

// Connect implements app.Storage.
func (s *Storage) Connect(_ context.Context) error {
	return nil
}

// CreateEvent implements app.Storage.
func (s *Storage) CreateEvent(event app.Event) (app.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	event.ID = uuid.New().String()
	if _, exists := s.events[event.ID]; !exists {
		// todo безсмысленная проверка, может надо в БЛ задавать ID ??
		return app.Event{}, app.ErrAlreadyExist
	}

	s.events[event.ID] = event
	return event, nil
}

// UpdateEvent implements app.Storage.
func (s *Storage) UpdateEvent(event app.Event) (app.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.events[event.ID]; !exists {
		return event, app.ErrNotFound
	}

	new_event := s.events[event.ID]
	new_event.ID = event.ID

	if event.Title != "" {
		new_event.Title = event.Title
	}
	if !event.Datetime.IsZero() {
		new_event.Datetime = event.Datetime
	}
	if event.Duration != 0 {
		new_event.Duration = event.Duration
	}
	if event.Description != "" {
		new_event.Description = event.Description
	}
	if event.UserID != "" {
		new_event.UserID = event.UserID
	}
	if event.NotifyTime != 0 {
		new_event.NotifyTime = event.NotifyTime
	}

	s.events[event.ID] = new_event
	return new_event, nil
}

// DeleteEvent implements app.Storage.
func (s *Storage) DeleteEvent(event app.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.events[event.ID]; !exists {
		// todo может не надо возвращать ошибку, ну нет события так нет ??
		return app.ErrNotFound
	}

	return nil
}

// ListEventsForDay implements app.Storage.
func (s *Storage) ListEventsForDay(date time.Time) []app.Event {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var events []app.Event
	for _, event := range s.events {
		isSame := event.Datetime.Year() == date.Year() &&
			event.Datetime.Month() == date.Month() &&
			event.Datetime.Day() == date.Day()

		if isSame {
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
func (s *Storage) ListEventsForMonth(date time.Time) []app.Event {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var events []app.Event
	weekStart := date.AddDate(0, 0, -int(date.Month()))
	weekEnd := weekStart.AddDate(0, 0, 7)

	for _, event := range s.events {
		if event.Datetime.After(weekStart) && event.Datetime.Before(weekEnd) {
			events = append(events, event)
		}
	}

	return events
}
