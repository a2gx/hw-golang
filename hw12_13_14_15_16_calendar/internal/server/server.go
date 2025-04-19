package server

import (
	"context"
	"fmt"
	"strings"

	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/internal/app"
	servergrpc "github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/internal/server/grpc"
	serverhttp "github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/internal/server/http"
	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/pkg/logger"
)

type Options struct {
	ServerType string
	HTTPAddr   string
	GRPCAddr   string
	Logg       *logger.Logger
	App        *app.App
}

type MultiServer struct {
	http   app.Server
	grpc   app.Server
	logg   *logger.Logger
	ctx    context.Context
	cancel context.CancelFunc
}

var _ app.Server = &MultiServer{}

func New(opts Options) (app.Server, error) {
	srv := &MultiServer{
		logg: opts.Logg,
	}

	switch strings.ToLower(opts.ServerType) {
	case "grpc":
		srv.grpc = servergrpc.New(opts.GRPCAddr, opts.Logg, opts.App)
	case "http":
		srv.http = serverhttp.New(opts.HTTPAddr, opts.Logg, opts.App)
	case "both":
		srv.grpc = servergrpc.New(opts.GRPCAddr, opts.Logg, opts.App)
		srv.http = serverhttp.New(opts.HTTPAddr, opts.Logg, opts.App)
	default:
		return nil, fmt.Errorf("unknown server type: %s", opts.ServerType)
	}

	return srv, nil
}

func (s *MultiServer) Start(ctx context.Context) error {
	s.ctx, s.cancel = context.WithCancel(ctx)
	errCh := make(chan error, 2)

	if s.http != nil {
		go func() {
			if err := s.http.Start(s.ctx); err != nil {
				errCh <- fmt.Errorf("failed to start http server: %w", err)
			}
		}()
	}
	if s.grpc != nil {
		go func() {
			if err := s.grpc.Start(s.ctx); err != nil {
				errCh <- fmt.Errorf("failed to start grpc server: %w", err)
			}
		}()
	}

	select {
	case err := <-errCh:
		s.cancel()
		return err
	case <-s.ctx.Done():
		return nil
	}
}

func (s *MultiServer) Stop(ctx context.Context) error {
	if s.cancel != nil {
		s.cancel()
	}

	if s.http != nil {
		if err := s.http.Stop(ctx); err != nil {
			s.logg.Error("failed to stop http server: " + err.Error())
		}
	}
	if s.grpc != nil {
		if err := s.grpc.Stop(ctx); err != nil {
			s.logg.Error("failed to stop grpc server: " + err.Error())
		}
	}

	return nil
}
