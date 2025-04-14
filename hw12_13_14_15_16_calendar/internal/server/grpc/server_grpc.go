package server_grpc

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"google.golang.org/grpc"

	pb "github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/api"
	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/internal/app"
	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/pkg/logger"
)

type Server struct {
	pb.UnimplementedCalendarServer
	grpcServer *grpc.Server
	logg       *logger.Logger
	app        *app.App
	addr       string
}

var _ app.Server = &Server{}

func New(addr string, logg *logger.Logger, app *app.App) *Server {
	return &Server{
		logg: logg,
		app:  app,
		addr: addr,
	}
}

func (s *Server) Start(_ context.Context) error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	s.grpcServer = grpc.NewServer()
	pb.RegisterCalendarServer(s.grpcServer, s)

	s.logg.Info("start GRPC server", slog.String("addr", s.addr))
	go func() {
		if err := s.grpcServer.Serve(listener); err != nil {
			s.logg.Error("failed GRPC to serve: %v", err)
		}
	}()

	return nil
}

func (s *Server) Stop(_ context.Context) error {
	if s.grpcServer != nil {
		s.logg.Info("GRPC server is stopping")
		s.grpcServer.GracefulStop()
	}

	return nil
}

func (s *Server) CreateEvent(ctx context.Context, req *pb.CreateEventRequest) (*pb.CreateEventReply, error) {
	s.logg.Info("handler CreateEvent")
	return &pb.CreateEventReply{}, nil
}
