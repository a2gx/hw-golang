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
	return a.store.CreateEvent(event)
}

func (a *App) UpdateEvent(event Event) (Event, error) {
	a.logg.Debug("App.UpdateEvent", "event", event)
	return a.store.UpdateEvent(event)
}

func (a *App) DeleteEvent(event Event) error {
	a.logg.Debug("App.DeleteEvent", "event", event)
	return a.store.DeleteEvent(event)
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
