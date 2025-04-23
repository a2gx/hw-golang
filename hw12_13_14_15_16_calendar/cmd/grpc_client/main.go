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
	createReq := &pb.CreateEventRequest{
		Title:       "Meeting",
		Description: "Team meeting to discuss project updates",
		StartTime:   time.Now().Format(time.RFC3339),
		EndTime:     time.Now().Add(1 * time.Hour).Format(time.RFC3339),
		NotifyTime:  time.Now().Add(-10 * time.Minute).Format(time.RFC3339),
	}
	createResp, err := c.CreateEvent(ctx, createReq)
	if err != nil {
		slog.Error("CreateEvent failed: " + err.Error())
		return
	}
	slog.Info("CreateEvent success: " + createResp.String())

	// UpdateEvent
	updateReq := &pb.UpdateEventRequest{
		Id:          createResp.Event.Id,
		Title:       "Updated Meeting",
		Description: "Updated description for the meeting",
		StartTime:   time.Now().Add(2 * time.Hour).Format(time.RFC3339),
		EndTime:     time.Now().Add(3 * time.Hour).Format(time.RFC3339),
		NotifyTime:  time.Now().Add(-15 * time.Minute).Format(time.RFC3339),
	}
	updateResp, err := c.UpdateEvent(ctx, updateReq)
	if err != nil {
		slog.Error("UpdateEvent failed: " + err.Error())
		return
	}
	slog.Info("UpdateEvent success: " + updateResp.String())
}
