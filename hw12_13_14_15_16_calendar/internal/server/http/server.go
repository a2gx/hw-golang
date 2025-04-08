package serverhttp

import (
	"context"
	"net/http"
	"time"

	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/internal/app"
	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/pkg/logger"
)

type Server struct {
	httpserver *http.Server
	logg       *logger.Logger
	app        *app.App
}

func New(logg *logger.Logger, app *app.App) *Server {
	mux := http.NewServeMux()
	hFn := &Handler{
		logg: logg,
		app:  app,
	}

	// register handlers ...
	mux.HandleFunc("/ping", hFn.GetPing)

	httpserver := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,

		Handler: applyMiddleware(
			mux,
			LoggerMiddleware(*logg),
		),
	}

	return &Server{
		httpserver: httpserver,
		logg:       logg,
		app:        app,
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.logg.Info("start http server")
	errCh := make(chan error)

	go func() {
		if err := s.httpserver.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		return s.Stop(context.Background())
	case err := <-errCh:
		return err
	}
}

func (s *Server) Stop(ctx context.Context) error {
	s.logg.Info("stop http server")
	return s.httpserver.Shutdown(ctx)
}
