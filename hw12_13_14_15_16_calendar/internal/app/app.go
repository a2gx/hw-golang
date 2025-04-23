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
	if updates.UserID != "" {
		existing.UserID = updates.UserID
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

func (a *App) getDateInterval(date time.Time, add int) (start, finish time.Time) {
	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	return date, date.AddDate(0, 0, add)
}

func (a *App) Create(event Event) (Event, error) {
	a.logg.Debug("App.Create", "event", event)
	return a.store.CreateEvent(event)
}

func (a *App) Update(event Event) (Event, error) {
	a.logg.Debug("App.Update", "event", event)

	if event.ID == "" {
		return Event{}, ErrIdRequired
	}

	existingEvent, err := a.store.GetById(event.ID)
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

func (a *App) ListEventsInInterval(from time.Time, days int) []Event {
	a.logg.Debug("App.ListEventsInInterval", "start_date", from, "days", days)

	st, fn := a.getDateInterval(from, days)
	return a.store.ListEventsInInterval(st, fn)
}

func (a *App) ListEventsForDay(date time.Time) []Event {
	a.logg.Debug("App.ListEventsForDay", "start_date", date)
	return a.store.ListEventsForDay(date)
}

func (a *App) ListEventsForWeek(date time.Time) []Event {
	a.logg.Debug("App.ListEventsForWeek", "start_date", date)
	return a.store.ListEventsForWeek(date)
}

func (a *App) ListEventsForMonth(date time.Time) []Event {
	a.logg.Debug("App.ListEventsForWeek", "start_date", date)
	return a.store.ListEventsForMonth(date)
}
