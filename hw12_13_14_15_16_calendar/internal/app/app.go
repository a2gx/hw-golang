package app

import (
	"time"

	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/pkg/logger"
)

type App struct {
	logg  *logger.Logger
	store Storage
}

var _ Application = &App{}

func New(logg *logger.Logger, store Storage) *App {
	return &App{
		logg:  logg,
		store: store,
	}
}

func (a *App) mergeEvents(existing, updates Event) Event {
	if updates.Title != "" {
		existing.Title = updates.Title
	}
	if updates.Description != "" {
		existing.Description = updates.Description
	}
	if !updates.StartTime.IsZero() {
		existing.StartTime = updates.StartTime
	}
	if !updates.EndTime.IsZero() {
		existing.EndTime = updates.EndTime
	}
	if !updates.NotifyTime.IsZero() {
		existing.NotifyTime = updates.NotifyTime
	}
	return existing
}

func (a *App) Create(event Event) (Event, error) {
	a.logg.Debug("App.Create", "event", event)
	return a.store.CreateEvent(event)
}

func (a *App) Update(event Event) (Event, error) {
	a.logg.Debug("App.Update", "event", event)

	if event.ID == "" {
		return Event{}, ErrIDRequired
	}

	existingEvent, err := a.store.GetByID(event.ID)
	if err != nil {
		return Event{}, err
	}

	event = a.mergeEvents(existingEvent, event)
	return a.store.UpdateEvent(event)
}

func (a *App) Delete(event Event) error {
	a.logg.Debug("App.Delete", "event", event)
	return a.store.DeleteEvent(event)
}

func (a *App) EventsInInterval(date time.Time, days int) []Event {
	a.logg.Debug("App.ListEventsInInterval", "start_date", date, "days", days)

	from := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	to := from.AddDate(0, 0, days)

	return a.store.FilterByInterval(from, to)
}
