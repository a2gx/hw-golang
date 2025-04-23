package serverhttp

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/internal/app"
	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/pkg/logger"
	"github.com/go-playground/validator/v10"
	"io"
	"net/http"
	"time"
)

type Handler struct {
	logg *logger.Logger
	app  *app.App
}

func (h *Handler) jsonReadAndValidate(body io.ReadCloser, input any) error {
	if err := json.NewDecoder(body).Decode(input); err != nil {
		return fmt.Errorf("invalid request body: %v", err)
	}
	if err := validator.New().Struct(input); err != nil {
		return fmt.Errorf("invalid validate: %v", err)
	}

	return nil
}

func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Title       string    `json:"title" validate:"required"`
		Description string    `json:"description,omitempty"`
		UserID      string    `json:"user_id" validate:"required"`
		StartTime   time.Time `json:"start_time" validate:"required"`
		EndTime     time.Time `json:"end_time" validate:"required,gtfield=StartTime"`
		NotifyTime  time.Time `json:"notify_time,omitempty"`
	}

	if err := h.jsonReadAndValidate(r.Body, &in); err != nil {
		h.logg.Error("Invalid JSON", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	event, err := h.app.Create(app.Event{
		Title:       in.Title,
		Description: in.Description,
		UserID:      in.UserID,
		StartTime:   in.StartTime,
		EndTime:     in.EndTime,
		NotifyTime:  in.NotifyTime,
	})
	if err != nil {
		h.logg.Error("failed to create event", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(event)
	if err != nil {
		h.logg.Error("failed to marshal response", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(res)
}

func (h *Handler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	var in struct {
		ID          string    `json:"id" validate:"required"`
		Title       string    `json:"title,omitempty"`
		Description string    `json:"description,omitempty" validate:"required"`
		UserID      string    `json:"user_id,omitempty"`
		StartTime   time.Time `json:"start_time,omitempty"`
		EndTime     time.Time `json:"end_time,omitempty"`
		NotifyTime  time.Time `json:"notify_time,omitempty"`
	}

	if err := h.jsonReadAndValidate(r.Body, &in); err != nil {
		h.logg.Error("Invalid JSON", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updEvent, err := h.app.Update(app.Event{
		ID:          in.ID,
		Title:       in.Title,
		Description: in.Description,
		UserID:      in.UserID,
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

	res, err := json.Marshal(updEvent)
	if err != nil {
		h.logg.Error("failed to marshal response", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(res)
}

func (h *Handler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	var in struct {
		ID string `json:"id" validate:"required"`
	}

	if err := h.jsonReadAndValidate(r.Body, &in); err != nil {
		h.logg.Error("Invalid JSON", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.app.Delete(app.Event{
		ID: in.ID,
	})
	if err != nil {
		h.logg.Error("failed to delete event", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) ListEventsForDay(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Date string `json:"date" validate:"required"`
	}

	if err := h.jsonReadAndValidate(r.Body, &in); err != nil {
		h.logg.Error("Invalid JSON", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	date, err := time.Parse("2006-01-02", in.Date)
	if err != nil {
		h.logg.Error("failed to parse date", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	events := h.app.ListEventsInInterval(date, 1)
	res, err := json.Marshal(events)
	if err != nil {
		h.logg.Error("failed to marshal response", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(res)
}

func (h *Handler) ListEventsForWeek(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Date string `json:"date" validate:"required"`
	}

	if err := h.jsonReadAndValidate(r.Body, &in); err != nil {
		h.logg.Error("Invalid JSON", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	date, err := time.Parse("2006-01-02", in.Date)
	if err != nil {
		h.logg.Error("failed to parse date", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	events := h.app.ListEventsInInterval(date, 7)
	res, err := json.Marshal(events)
	if err != nil {
		h.logg.Error("failed to marshal response", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(res)
}

func (h *Handler) ListEventsForMonth(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Date string `json:"date" validate:"required"`
	}

	if err := h.jsonReadAndValidate(r.Body, &in); err != nil {
		h.logg.Error("Invalid JSON", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	date, err := time.Parse("2006-01-02", in.Date)
	if err != nil {
		h.logg.Error("failed to parse date", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	events := h.app.ListEventsInInterval(date, 7)
	res, err := json.Marshal(events)
	if err != nil {
		h.logg.Error("failed to marshal response", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(res)
}
