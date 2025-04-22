package serverhttp

import (
	"net/http"

	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/internal/app"
	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/pkg/logger"
)

type Handler struct {
	logg *logger.Logger
	app  *app.App
}

func (h *Handler) Home(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("home"))
}

func (h *Handler) Ping(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("pong"))
}

func (h *Handler) CreateEvent(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("CreateEvent"))
}

func (h *Handler) UpdateEvent(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("UpdateEvent"))
}

func (h *Handler) DeleteEvent(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("DeleteEvent"))
}

func (h *Handler) ListEventsForDay(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ListEventsForDay"))
}

func (h *Handler) ListEventsForWeek(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ListEventsForWeek"))
}

func (h *Handler) ListEventsForMonth(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ListEventsForMonth"))
}
