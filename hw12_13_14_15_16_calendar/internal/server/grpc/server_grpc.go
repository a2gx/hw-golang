package servergrpc

import (
	"context"
	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/internal/app"
	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/pkg/logger"
)

type Server struct{}

var _ app.Server = &Server{}

func New(addr string, logg *logger.Logger, app *app.App) *Server {
	return &Server{}
}

func (s Server) Start(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (s Server) Stop(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}
