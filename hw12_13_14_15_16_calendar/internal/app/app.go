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

func (a *App) CreateEvent(event Event) (Event, error) {
	a.logg.Debug("App.CreateEvent", "event", event)
	return event, nil
}

func (a *App) UpdateEvent(event Event) (Event, error) {
	a.logg.Debug("App.UpdateEvent", "event", event)
	return event, nil
}

func (a *App) DeleteEvent(event Event) error {
	a.logg.Debug("App.DeleteEvent", "event", event)
	return nil
}

func (a *App) ListEventsForDay(date time.Time) []Event {
	a.logg.Debug("App.ListEventsForDay", "start_date", date)
	events := make([]Event, 0)
	return events
}

func (a *App) ListEventsForWeek(date time.Time) []Event {
	a.logg.Debug("App.ListEventsForWeek", "start_date", date)
	events := make([]Event, 0)
	return events
}

func (a *App) ListEventsForMonth(date time.Time) []Event {
	a.logg.Debug("App.ListEventsForWeek", "start_date", date)
	events := make([]Event, 0)
	return events
}
