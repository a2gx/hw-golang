package servergrpc

import (
	"context"
	"time"

	pb "github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/api"
	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/internal/app"
	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/pkg/logger"
)

type Handler struct {
	pb.UnimplementedCalendarServer
	logg *logger.Logger
	app  *app.App
}

func (h *Handler) CreateEvent(_ context.Context, req *pb.CreateEventRequest) (*pb.CreateEventReply, error) {
	h.logg.Info("gRPC CreateEvent called", "request", req)

	event, err := h.app.Create(app.Event{
		Title:       req.Title,
		Description: req.Description,
		StartTime:   parseTime(req.StartTime),
		EndTime:     parseTime(req.EndTime),
		NotifyTime:  parseTime(req.NotifyTime),
	})
	if err != nil {
		h.logg.Error("failed to create event", "error", err)
		return nil, err
	}

	return &pb.CreateEventReply{
		Event: toProtoEvent(event),
	}, nil
}

func (h *Handler) UpdateEvent(_ context.Context, req *pb.UpdateEventRequest) (*pb.UpdateEventReply, error) {
	h.logg.Info("gRPC UpdateEvent called", "request", req)

	event, err := h.app.Update(app.Event{
		ID:          req.Id,
		Title:       req.Title,
		Description: req.Description,
		StartTime:   parseTime(req.StartTime),
		EndTime:     parseTime(req.EndTime),
		NotifyTime:  parseTime(req.NotifyTime),
	})
	if err != nil {
		h.logg.Error("failed to update event", "error", err)
		return nil, err
	}

	return &pb.UpdateEventReply{
		Event: toProtoEvent(event),
	}, nil
}

func (h *Handler) DeleteEvent(_ context.Context, req *pb.DeleteEventRequest) (*pb.DeleteEventReply, error) {
	h.logg.Info("gRPC DeleteEvent called", "request", req)

	err := h.app.Delete(app.Event{ID: req.Id})
	if err != nil {
		h.logg.Error("failed to delete event", "error", err)
		return nil, err
	}

	return &pb.DeleteEventReply{}, nil
}

func (h *Handler) ListEventsForDay(
	_ context.Context,
	req *pb.ListEventsForDayRequest,
) (*pb.ListEventsForDayReply, error) {
	h.logg.Info("gRPC ListEventsForDay called", "request", req)

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		h.logg.Error("invalid date format", "error", err)
		return nil, err
	}

	events := h.app.EventsInInterval(date, 1)
	return &pb.ListEventsForDayReply{
		Events: toProtoEvents(events),
	}, nil
}

func (h *Handler) ListEventsForWeek(
	_ context.Context,
	req *pb.ListEventsForWeekRequest,
) (*pb.ListEventsForWeekReply, error) {
	h.logg.Info("gRPC ListEventsForWeek called", "request", req)

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		h.logg.Error("invalid date format", "error", err)
		return nil, err
	}

	events := h.app.EventsInInterval(date, 7)
	return &pb.ListEventsForWeekReply{
		Events: toProtoEvents(events),
	}, nil
}

func (h *Handler) ListEventsForMonth(
	_ context.Context,
	req *pb.ListEventsForMonthRequest,
) (*pb.ListEventsForMonthReply, error) {
	h.logg.Info("gRPC ListEventsForMonth called", "request", req)

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		h.logg.Error("invalid date format", "error", err)
		return nil, err
	}

	events := h.app.EventsInInterval(date, 30)
	return &pb.ListEventsForMonthReply{
		Events: toProtoEvents(events),
	}, nil
}

func parseTime(t string) time.Time {
	parsed, _ := time.Parse(time.RFC3339, t)
	return parsed
}

func toProtoEvent(event app.Event) *pb.Event {
	return &pb.Event{
		Id:          event.ID,
		Title:       event.Title,
		Description: event.Description,
		StartTime:   event.StartTime.Format(time.RFC3339),
		EndTime:     event.EndTime.Format(time.RFC3339),
		NotifyTime:  event.NotifyTime.Format(time.RFC3339),
	}
}

func toProtoEvents(events []app.Event) []*pb.Event {
	protoEvents := make([]*pb.Event, len(events))
	for i, event := range events {
		protoEvents[i] = toProtoEvent(event)
	}
	return protoEvents
}
