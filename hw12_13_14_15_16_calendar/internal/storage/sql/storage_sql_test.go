package storagesql

import (
	_ "database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/internal/app"
	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/pkg/logger"
	"github.com/stretchr/testify/require"
)

func TestStorage_CreateEvent(t *testing.T) {
	logg := logger.New(logger.Options{Level: "debug"})
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	storage := &Storage{logg: logg, db: db}

	event := app.Event{
		Title:       "Test Event",
		Description: "Test Description",
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(1 * time.Hour),
		NotifyTime:  time.Now().Add(-10 * time.Minute),
	}

	mock.ExpectQuery("INSERT INTO events").
		WithArgs(event.Title, event.Description, event.StartTime, event.EndTime, event.NotifyTime).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("1"))

	createdEvent, err := storage.CreateEvent(event)
	require.NoError(t, err)
	require.Equal(t, "1", createdEvent.ID)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStorage_GetByID(t *testing.T) {
	logg := logger.New(logger.Options{Level: "debug"})
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	storage := &Storage{logg: logg, db: db}

	eventID := "1"
	mock.ExpectQuery("SELECT id, title, description, start_time, end_time FROM events WHERE id = \\$1").
		WithArgs(eventID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "start_time", "end_time"}).
			AddRow(eventID, "Test Event", "Test Description", time.Now(), time.Now().Add(1*time.Hour)))

	event, err := storage.GetByID(eventID)
	require.NoError(t, err)
	require.Equal(t, eventID, event.ID)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStorage_DeleteEvent(t *testing.T) {
	logg := logger.New(logger.Options{Level: "debug"})
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	storage := &Storage{logg: logg, db: db}

	event := app.Event{ID: "1"}
	mock.ExpectExec("DELETE FROM events WHERE id = \\$1").
		WithArgs(event.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = storage.DeleteEvent(event)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStorage_FetchEventsToNotify(t *testing.T) {
	logg := logger.New(logger.Options{Level: "debug"})
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	storage := &Storage{logg: logg, db: db}

	mock.ExpectQuery("SELECT id, title, description, start_time, end_time, notify_time FROM events WHERE notify_time <= NOW()").
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "start_time", "end_time", "notify_time"}).
			AddRow("1", "Test Event", "Test Description", time.Now(), time.Now().Add(1*time.Hour), time.Now().Add(-10*time.Minute)))

	events, err := storage.FetchEventsToNotify()
	require.NoError(t, err)
	require.Len(t, events, 1)
	require.Equal(t, "1", events[0].ID)
	require.NoError(t, mock.ExpectationsWereMet())
}
