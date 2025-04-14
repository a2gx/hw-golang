package server

import (
	"context"
	"fmt"

	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/internal/app"
	servergrpc "github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/internal/server/grpc"
	serverhttp "github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/internal/server/http"
	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/pkg/logger"
)

type Options struct {
	ServerType string
	HttpAddr   string
	GrpcAddr   string
	Logg       *logger.Logger
	App        *app.App
}

type MultiServer struct {
	httpServer app.Server
	grpcServer app.Server
	logg       *logger.Logger
}

var _ app.Server = &MultiServer{}

func New(opts Options) (app.Server, error) {
	srv := &MultiServer{
		logg: opts.Logg,
	}

	switch opts.ServerType {
	case "http":
		srv.httpServer = serverhttp.New(opts.Logg, opts.App, opts.HttpAddr)
	case "grpc":
		srv.grpcServer = servergrpc.New(opts.Logg, opts.App, opts.GrpcAddr)
	case "both":
		srv.httpServer = serverhttp.New(opts.Logg, opts.App, opts.HttpAddr)
		srv.grpcServer = servergrpc.New(opts.Logg, opts.App, opts.GrpcAddr)
	}

	return srv, nil
}

func (s *MultiServer) Start(ctx context.Context) error {
	errChan := make(chan error, 2)

	if s.httpServer != nil {
		go func() {
			if err := s.httpServer.Start(ctx); err != nil {
				errChan <- fmt.Errorf("failed to start http server: %w", err)
			}
		}()
	}

	if s.grpcServer != nil {
		go func() {
			if err := s.grpcServer.Start(ctx); err != nil {
				errChan <- fmt.Errorf("failed to start grpc server: %w", err)
			}
		}()
	}

	// Ждем либо ошибку, либо завершение контекста
	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (s *MultiServer) Stop(ctx context.Context) error {
	if s.httpServer != nil {
		if err := s.httpServer.Stop(ctx); err != nil {
			s.logg.Error("failed to stop http server: " + err.Error())
		}
	}

	if s.grpcServer != nil {
		if err := s.grpcServer.Stop(ctx); err != nil {
			s.logg.Error("failed to stop grpc server: " + err.Error())
		}
	}

	return nil
}
