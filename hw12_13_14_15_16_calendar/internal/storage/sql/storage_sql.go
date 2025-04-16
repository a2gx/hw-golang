package storagesql

import (
	"database/sql"
	"time"

	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/internal/app"
	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/pkg/logger"
	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/pkg/tools"
	_ "github.com/lib/pq" // Регистрируем драйвер PostgreSQL
)

const (
	querySelectEventsInInterval = `
		SELECT id, title, start_time, end_time, description 
		FROM events 
		WHERE start_time < $1 AND end_time > $2
	`
)

type Storage struct {
	logg *logger.Logger
	db   *sql.DB
	dns  string // Data Source Name
}

var _ app.Storage = &Storage{}

func New(logg *logger.Logger, dns string) *Storage {
	return &Storage{
		logg: logg,
		dns:  dns,
	}
}

func (s *Storage) Connect() error {
	db, err := sql.Open("postgres", s.dns)
	if err != nil {
		s.logg.Error("failed to connect to database", "error", err)
		return err
	}

	if err := db.Ping(); err != nil {
		s.logg.Error("failed to ping database", "error", err)
		return err
	}

	s.db = db
	s.logg.Debug("storage SQL connected")
	return nil
}

func (s *Storage) Close() error {
	if s.db == nil {
		return nil
	}

	if err := s.db.Close(); err != nil {
		s.logg.Error("failed to close database connection", "error", err)
		return err
	}

	s.logg.Debug("storage SQL closed")
	return nil
}

func (s *Storage) CreateEvent(event app.Event) (app.Event, error) {
	query := `INSERT INTO events (id, title, start_time, end_time, description) 
			  VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := s.db.QueryRow(query, event.ID, event.Title, event.StartTime, event.EndTime, event.Description).Scan(&event.ID)
	if err != nil {
		s.logg.Error("failed to create event", "error", err)
		return app.Event{}, err
	}

	s.logg.Debug("event created", "id", event.ID)
	return event, nil
}

func (s *Storage) UpdateEvent(event app.Event) (app.Event, error) {
	query := `UPDATE events SET title = $1, start_time = $2, end_time = $3, description = $4 WHERE id = $5`
	_, err := s.db.Exec(query, event.Title, event.StartTime, event.EndTime, event.Description, event.ID)
	if err != nil {
		s.logg.Error("failed to update event", "error", err)
		return app.Event{}, err
	}

	s.logg.Debug("event updated", "id", event.ID)
	return event, nil
}

func (s *Storage) DeleteEvent(event app.Event) error {
	query := `DELETE FROM events WHERE id = $1`
	_, err := s.db.Exec(query, event.ID)
	if err != nil {
		s.logg.Error("failed to delete event", "error", err)
		return err
	}

	s.logg.Debug("event deleted", "id", event.ID)
	return nil
}

func (s *Storage) ListEventsForDay(day time.Time) []app.Event {
	start, end := tools.GetDateInterval(day, 1)
	events, err := s.selectEventsInInterval(start, end)
	if err != nil {
		s.logg.Error("failed to list events for day", "error", err)
		return nil
	}

	s.logg.Debug("events listed for day", "start_date", start, "end_date", end, "count", len(events))
	return events
}

func (s *Storage) ListEventsForWeek(week time.Time) []app.Event {
	start, end := tools.GetDateInterval(week, 7)
	events, err := s.selectEventsInInterval(start, end)
	if err != nil {
		s.logg.Error("failed to list events for week", "error", err)
		return nil
	}

	s.logg.Debug("events listed for week", "start_date", start, "end_date", end, "count", len(events))
	return events
}

func (s *Storage) ListEventsForMonth(month time.Time) []app.Event {
	start, end := tools.GetDateInterval(month, 30)
	events, err := s.selectEventsInInterval(start, end)
	if err != nil {
		s.logg.Error("failed to list events for month", "error", err)
		return nil
	}

	s.logg.Debug("events listed for month", "start_date", start, "end_date", end, "count", len(events))
	return events
}

func (s *Storage) selectEventsInInterval(start, end time.Time) ([]app.Event, error) {
	query := querySelectEventsInInterval
	rows, err := s.db.Query(query, end, start)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []app.Event
	for rows.Next() {
		var event app.Event
		if err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.StartTime,
			&event.EndTime,
			&event.Description,
		); err != nil {
			s.logg.Error("failed to scan event", "error", err)
			continue
		}
		events = append(events, event)
	}

	return events, nil
}
