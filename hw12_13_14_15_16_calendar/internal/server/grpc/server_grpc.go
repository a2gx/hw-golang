package servergrpc

import (
	"context"

	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/internal/app"
	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/pkg/logger"
)

type Server struct {
	logg *logger.Logger
	app  *app.App
}

var _ app.Server = &Server{}

func New(addr string, logg *logger.Logger, app *app.App) *Server {
	return &Server{
		logg: logg,
		app:  app,
	}
}

func (s Server) Start(ctx context.Context) error {
	return nil
}

func (s Server) Stop(ctx context.Context) error {
	return nil
}
