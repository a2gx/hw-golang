package serverhttp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/internal/app"
	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/pkg/logger"
)

type Server struct {
	srv  *http.Server
	logg *logger.Logger
	app  *app.App
}

var _ app.Server = &Server{}

func New(addr string, logg *logger.Logger, app *app.App) *Server {
	mux := http.NewServeMux()
	h := &Handler{logg, app}

	mux.HandleFunc("POST /event", h.CreateEvent)
	mux.HandleFunc("PUT /event", h.UpdateEvent)
	mux.HandleFunc("DELETE /event", h.DeleteEvent)
	mux.HandleFunc("GET /events-day", h.ListEventsForDay)
	mux.HandleFunc("GET /events-week", h.ListEventsForWeek)
	mux.HandleFunc("GET /events-month", h.ListEventsForMonth)

	srv := &http.Server{
		Addr:         addr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,

		Handler: applyMiddleware(
			mux,
			loggingMiddleware,
			maxBytesMiddleware,
			panicMiddleware,
		),
	}

	return &Server{srv, logg, app}
}

func (s *Server) Start(ctx context.Context) error {
	s.logg.Debug("start HTTP server", "addr", s.srv.Addr)
	errCh := make(chan error)

	go func() {
		if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return fmt.Errorf("failed HTTP to serve: %w", err)
	case <-ctx.Done():
		return nil
	}
}

func (s *Server) Stop(ctx context.Context) error {
	s.logg.Debug("stop HTTP server")

	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := s.srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown HTTP server error: %w", err)
	}

	return nil
}
