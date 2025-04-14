package main

import (
	"context"
	"flag"
	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/internal/app"
	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/internal/server"
	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/internal/storage"
	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/pkg/logger"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	flag.Parse()
	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	// Initialize Configuration
	cfg, err := NewConfig()
	if err != nil {
		log.Fatalf("failed to init configuration: %v", err)
	}

	// Initialize Logger
	logg := logger.New(logger.Options{
		Level:       cfg.Logger.Level,
		HandlerType: cfg.Logger.HandlerType,
		Filename:    cfg.Logger.Filename,
		AddSource:   cfg.Logger.AddSource,
	})

	// Initialize Storage
	store, err := storage.New(storage.Options{
		StorageType: cfg.App.Storage,
		Logg:        logg,
	})
	if err != nil {
		log.Fatalf("failed to init storage: %v", err)
	}

	// Initialize Business-Logic
	calendar := app.New(logg, store)

	// Initialize HttpServer
	srv, err := server.New(server.Options{
		ServerType: cfg.App.Server,
		HttpAddr:   cfg.Server.HttpAddr,
		GrpcAddr:   cfg.Server.GrpcAddr,
		Logg:       logg,
		App:        calendar,
	})
	if err != nil {
		log.Fatalf("failed to init server: %v", err)
	}

	// Start server ...

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	logg.Info("calendar is running...")

	if err := srv.Start(context.Background()); err != nil {
		logg.Error("failed to start server: " + err.Error())
		os.Exit(1)
	}

	<-ctx.Done()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Second*3)
	defer shutdownCancel()

	if err := srv.Stop(shutdownCtx); err != nil {
		logg.Error("failed to stop server: " + err.Error())
	}
}
