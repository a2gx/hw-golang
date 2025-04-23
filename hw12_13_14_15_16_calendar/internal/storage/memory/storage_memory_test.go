package storagememory

import (
	"testing"
	"time"

	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/internal/app"
	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/pkg/logger"
)

func TestStorageMemory(t *testing.T) {
	logg := logger.New(logger.Options{Level: "info"})
	store := New(logg)

	if err := store.Connect(); err != nil {
		t.Fatalf("failed to connect to storage: %v", err)
	}
	defer store.Close()

	event := app.Event{
		Title:       "Event title",
		Description: "Event description",
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(1 * time.Hour),
	}

	createEvent, err := store.CreateEvent(event)
	if err != nil {
		t.Fatalf("failed to create event: %v", err)
	}
	if createEvent.ID == "" {
		t.Error("expected non-empty event ID")
	}

	events := store.FilterByInterval(time.Now(), time.Now().Add(1*time.Hour))
	if len(events) == 0 {
		t.Error("expected to find the created event")
	}

	createEvent.Title = "__UPDATED_TEST_EVENT__"
	updateEvent, err := store.UpdateEvent(createEvent)
	if err != nil {
		t.Fatalf("failed to update event: %v", err)
	}
	if updateEvent.Title != "__UPDATED_TEST_EVENT__" {
		t.Errorf("expected updated title to be 'Updated Test Event', got '%s'", updateEvent.Title)
	}

	err = store.DeleteEvent(createEvent)
	if err != nil {
		t.Fatalf("failed to delete event: %v", err)
	}

	events = store.FilterByInterval(time.Now(), time.Now().Add(1*time.Hour))
	if len(events) != 0 {
		t.Error("expected no events after deletion")
	}
}
