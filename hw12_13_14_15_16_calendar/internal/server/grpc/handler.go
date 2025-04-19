package servergrpc

import (
	"context"

	pb "github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/api"
	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/internal/app"
	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/pkg/logger"
)

type Handler struct {
	pb.UnimplementedCalendarServer
	logg *logger.Logger
	app  *app.App
}

func (h *Handler) CreateEvent(_ context.Context, _ *pb.CreateEventRequest) (*pb.CreateEventReply, error) {
	h.logg.Info("handler CreateEvent")
	return &pb.CreateEventReply{}, nil
}
