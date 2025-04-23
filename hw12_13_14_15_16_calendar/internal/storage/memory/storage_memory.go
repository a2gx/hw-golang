package storagememory

import (
	"sync"
	"time"

	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/internal/app"
	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/pkg/logger"
	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/pkg/tools"
	"github.com/google/uuid"
)

type Storage struct {
	mu     sync.RWMutex
	logg   *logger.Logger
	events map[string]app.Event
}

var _ app.Storage = &Storage{}

func New(logg *logger.Logger) *Storage {
	return &Storage{
		logg:   logg,
		events: make(map[string]app.Event),
	}
}

func (s *Storage) Connect() error {
	s.logg.Debug("storage Memory connected")
	return nil
}

func (s *Storage) Close() error {
	s.logg.Debug("storage Memory closed")
	return nil
}

func (s *Storage) GetById(eventId string) (app.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if event, exists := s.events[eventId]; !exists {
		return app.Event{}, app.ErrNotFound
	} else {
		return event, nil
	}
}

func (s *Storage) CreateEvent(event app.Event) (app.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	event.ID = uuid.New().String()
	s.events[event.ID] = event

	s.logg.Debug("event created", "id", event.ID)
	return event, nil
}

func (s *Storage) UpdateEvent(event app.Event) (app.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.events[event.ID]; !exists {
		return app.Event{}, app.ErrNotFound
	}
	s.events[event.ID] = event

	s.logg.Debug("event updated", "id", event.ID)
	return event, nil
}

func (s *Storage) DeleteEvent(event app.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.events, event.ID)

	s.logg.Debug("event deleted", "id", event.ID)
	return nil
}

func (s *Storage) ListEventsInInterval(st, fn time.Time) []app.Event {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []app.Event
	for _, e := range s.events {
		if e.StartTime.Before(fn) && e.EndTime.After(st) {
			result = append(result, e)
		}
	}

	s.logg.Debug("events listed for day", "start_date", st, "end_date", fn, "count", len(result))
	return result
}

// next need kill...

func (s *Storage) ListEventsForDay(date time.Time) []app.Event {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []app.Event
	st, fn := tools.GetDateInterval(date, 1) // +1 день, только указанный день

	for _, e := range s.events {
		if e.StartTime.Before(fn) && e.EndTime.After(st) {
			result = append(result, e)
		}
	}

	s.logg.Debug("events listed for day", "start_date", st, "end_date", fn, "count", len(result))
	return result
}

func (s *Storage) ListEventsForWeek(date time.Time) []app.Event {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []app.Event
	st, fn := tools.GetDateInterval(date, 7) // +7 дней

	for _, e := range s.events {
		if e.StartTime.Before(fn) && e.EndTime.After(st) {
			result = append(result, e)
		}
	}

	s.logg.Debug("events listed for week", "start_date", st, "end_date", fn, "count", len(result))
	return result
}

func (s *Storage) ListEventsForMonth(date time.Time) []app.Event {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []app.Event
	st, fn := tools.GetDateInterval(date, 30) // +30 дней

	for _, e := range s.events {
		if e.StartTime.Before(fn) && e.EndTime.After(st) {
			result = append(result, e)
		}
	}

	s.logg.Debug("events listed for month", "start_date", st, "end_date", fn, "count", len(result))
	return result
}
