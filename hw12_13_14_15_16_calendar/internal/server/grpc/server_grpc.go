package servergrpc

import (
	"context"
	"fmt"
	"net"

	pb "github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/api"
	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/internal/app"
	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/pkg/logger"
	"google.golang.org/grpc"
)

type Server struct {
	grpcServer *grpc.Server
	logg       *logger.Logger
	handler    *Handler
	addr       string
}

var _ app.Server = &Server{}

func New(addr string, logg *logger.Logger, app *app.App) *Server {
	return &Server{
		logg: logg,
		addr: addr,
		handler: &Handler{
			logg: logg,
			app:  app,
		},
	}
}

func (s *Server) Start(ctx context.Context) error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("failed GRPC to listen: %w", err)
	}

	s.grpcServer = grpc.NewServer(
		grpc.UnaryInterceptor(loggingInterceptor(s.logg)),
	)
	pb.RegisterCalendarServer(s.grpcServer, s.handler)

	s.logg.Debug("start GRPC server", "addr", s.addr)
	errCh := make(chan error)

	go func() {
		if err := s.grpcServer.Serve(listener); err != nil {
			errCh <- fmt.Errorf("failed GRPC to serve: %w", err)
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		return nil
	}
}

func (s *Server) Stop(ctx context.Context) error {
	s.logg.Debug("stop GRPC server")
	stopped := make(chan struct{})

	go func() {
		s.grpcServer.GracefulStop()
		close(stopped)
	}()

	select {
	case <-ctx.Done():
		s.grpcServer.Stop()
	case <-stopped:
	}

	return nil
}
