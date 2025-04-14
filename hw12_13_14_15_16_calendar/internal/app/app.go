package app

import (
	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/pkg/logger"
	"time"
)

type App struct {
	logg  *logger.Logger
	store Storage
}

var _ Application = &App{} //

func New(logg *logger.Logger, store Storage) *App {
	return &App{
		logg:  logg,
		store: store,
	}
}

func (a *App) CreateEvent(event Event) (Event, error) {
	//TODO implement me
	panic("implement me")
}

func (a *App) UpdateEvent(event Event) (Event, error) {
	//TODO implement me
	panic("implement me")
}

func (a *App) DeleteEvent(event Event) error {
	//TODO implement me
	panic("implement me")
}

func (a *App) ListEventsForDay(day time.Time) []Event {
	//TODO implement me
	panic("implement me")
}

func (a *App) ListEventsForWeek(week time.Time) []Event {
	//TODO implement me
	panic("implement me")
}

func (a *App) ListEventsForMonth(month time.Time) []Event {
	//TODO implement me
	panic("implement me")
}
