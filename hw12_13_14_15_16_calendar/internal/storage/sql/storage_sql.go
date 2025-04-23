package storagesql

import (
	"database/sql"
	"time"

	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/internal/app"
	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/pkg/logger"
	_ "github.com/lib/pq" // Регистрируем драйвер PostgreSQL
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
	query := `
		INSERT INTO events (title, description, start_time, end_time, notify_time) 
		VALUES ($1, $2, $3, $4, $5) RETURNING id
`
	err := s.db.QueryRow(
		query,
		event.Title,
		event.Description,
		event.StartTime,
		event.EndTime,
		event.NotifyTime,
	).Scan(&event.ID)
	if err != nil {
		s.logg.Error("failed to create event", "error", err)
		return app.Event{}, err
	}

	s.logg.Debug("event created", "id", event.ID)
	return event, nil
}

func (s *Storage) UpdateEvent(event app.Event) (app.Event, error) {
	query := `
		UPDATE events SET title = $1, start_time = $2, end_time = $3, description = $4 
		WHERE id = $5
	`
	_, err := s.db.Exec(query, event.Title, event.StartTime, event.EndTime, event.Description, event.ID)
	if err != nil {
		s.logg.Error("failed to update event", "error", err)
		return app.Event{}, err
	}

	s.logg.Debug("event updated", "id", event.ID)
	return event, nil
}

func (s *Storage) DeleteEvent(event app.Event) error {
	query := `
		DELETE FROM events 
		WHERE id = $1
	`
	_, err := s.db.Exec(query, event.ID)
	if err != nil {
		s.logg.Error("failed to delete event", "error", err)
		return err
	}

	s.logg.Debug("event deleted", "id", event.ID)
	return nil
}

func (s *Storage) GetByID(eventID string) (app.Event, error) {
	query := `
		SELECT id, title, description, start_time, end_time 
		FROM events 
		WHERE id = $1
	`
	row := s.db.QueryRow(query, eventID)

	var event app.Event
	if err := row.Scan(
		&event.ID,
		&event.Title,
		&event.Description,
		&event.StartTime,
		&event.EndTime,
	); err != nil {
		s.logg.Error("failed to get event by ID", "error", err)
		return app.Event{}, err
	}

	s.logg.Debug("event retrieved", "id", event.ID)
	return event, nil
}

func (s *Storage) FilterByInterval(start, end time.Time) []app.Event {
	query := `
		SELECT id, title, description, start_time, end_time 
		FROM events 
		WHERE start_time < $1 AND end_time > $2
	`
	rows, err := s.db.Query(query, end, start)
	if err != nil {
		s.logg.Error("failed to list events", "error", err)
		return nil
	}
	defer rows.Close()

	var events []app.Event
	for rows.Next() {
		var event app.Event
		if err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.Description,
			&event.StartTime,
			&event.EndTime,
		); err != nil {
			s.logg.Error("failed to scan event", "error", err)
			continue
		}
		events = append(events, event)
	}

	s.logg.Debug("events listed for day", "start_date", start, "end_date", end, "count", len(events))
	return events
}

func (s *Storage) FetchEventsToNotify() ([]app.Event, error) {
	query := `
		SELECT id, title, description, start_time, end_time, notify_time
		FROM events
		WHERE notify_time <= NOW()
	`
	rows, err := s.db.Query(query)
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
			&event.Description,
			&event.StartTime,
			&event.EndTime,
			&event.NotifyTime,
		); err != nil {
			s.logg.Error("failed to scan event", "error", err)
			continue
		}
		events = append(events, event)
	}

	return events, nil
}

func (s *Storage) DeleteOldEvents() error {
	query := `DELETE FROM events
		WHERE start_time < NOW() - INTERVAL '1 year'`
	_, err := s.db.Exec(query)
	if err != nil {
		s.logg.Error("failed to delete old events", "error", err)
		return err
	}

	return nil
}
