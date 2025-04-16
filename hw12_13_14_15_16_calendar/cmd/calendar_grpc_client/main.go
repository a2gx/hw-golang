package main

import (
	"context"
	"log"
	"log/slog"
	"time"

	pb "github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg, err := NewConfig()
	if err != nil {
		log.Fatalf("failed to init configuration: %v", err)
	}

	conn, err := grpc.NewClient(cfg.Server.GRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewCalendarClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// CreateEvent
	if r, err := c.CreateEvent(ctx, &pb.CreateEventRequest{}); err != nil {
		slog.Error("CreateEvent failed: " + err.Error())
	} else {
		slog.Info("CreateEvent success: " + r.String())
	}
}
