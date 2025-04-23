package serverhttp

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/internal/app"
	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/pkg/logger"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	logg *logger.Logger
	app  *app.App
}

func (h *Handler) jsonResponse(w http.ResponseWriter, data any) {
	w.WriteHeader(http.StatusOK)

	if data == nil {
		return
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logg.Error("failed to write response", "error", err)
	}
}

func (h *Handler) jsonReadAndValidate(r *http.Request, input any) error {
	if err := json.NewDecoder(r.Body).Decode(input); err != nil {
		return fmt.Errorf("invalid request body: %w", err)
	}
	if err := validator.New().Struct(input); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	return nil
}

func (h *Handler) handleEventsInInterval(w http.ResponseWriter, r *http.Request, days int) {
	var in struct {
		Date string `json:"date" validate:"required"`
	}

	if err := h.jsonReadAndValidate(r, &in); err != nil {
		h.logg.Error("request error", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	date, err := time.Parse("2006-01-02", in.Date)
	if err != nil {
		h.logg.Error("invalid date format", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	events := h.app.EventsInInterval(date, days)
	h.jsonResponse(w, events)
}

func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Title       string    `json:"title" validate:"required"`
		Description string    `json:"description,omitempty"`
		StartTime   time.Time `json:"start_time" validate:"required"`
		EndTime     time.Time `json:"end_time" validate:"required,gtfield=StartTime"`
		NotifyTime  time.Time `json:"notify_time,omitempty"`
	}

	if err := h.jsonReadAndValidate(r, &in); err != nil {
		h.logg.Error("request error", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	event, err := h.app.Create(app.Event{
		Title:       in.Title,
		Description: in.Description,
		StartTime:   in.StartTime,
		EndTime:     in.EndTime,
		NotifyTime:  in.NotifyTime,
	})
	if err != nil {
		h.logg.Error("failed to create event", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.jsonResponse(w, event)
}

func (h *Handler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	var in struct {
		ID          string    `json:"id" validate:"required"`
		Title       string    `json:"title,omitempty"`
		Description string    `json:"description,omitempty"`
		StartTime   time.Time `json:"start_time,omitempty"`
		EndTime     time.Time `json:"end_time,omitempty"`
		NotifyTime  time.Time `json:"notify_time,omitempty"`
	}

	if err := h.jsonReadAndValidate(r, &in); err != nil {
		h.logg.Error("request error", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	event, err := h.app.Update(app.Event{
		ID:          in.ID,
		Title:       in.Title,
		Description: in.Description,
		StartTime:   in.StartTime,
		EndTime:     in.EndTime,
		NotifyTime:  in.NotifyTime,
	})
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, app.ErrNotFound) {
			status = http.StatusNotFound
		}
		h.logg.Error("failed to update event", "error", err)
		http.Error(w, err.Error(), status)
		return
	}

	h.jsonResponse(w, event)
}

func (h *Handler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	var in struct {
		ID string `json:"id" validate:"required"`
	}

	if err := h.jsonReadAndValidate(r, &in); err != nil {
		h.logg.Error("request error", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.app.Delete(app.Event{ID: in.ID}); err != nil {
		h.logg.Error("failed to delete event", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.jsonResponse(w, nil)
}

func (h *Handler) ListEventsForDay(w http.ResponseWriter, r *http.Request) {
	h.handleEventsInInterval(w, r, 1)
}

func (h *Handler) ListEventsForWeek(w http.ResponseWriter, r *http.Request) {
	h.handleEventsInInterval(w, r, 7)
}

func (h *Handler) ListEventsForMonth(w http.ResponseWriter, r *http.Request) {
	h.handleEventsInInterval(w, r, 30)
}
