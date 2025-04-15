package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/internal/app"
	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/internal/server"
	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/internal/storage"
	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/pkg/logger"
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
		Level:     cfg.Logger.Level,
		Handler:   cfg.Logger.Handler,
		Filename:  cfg.Logger.Filename,
		AddSource: cfg.Logger.AddSource,
	})

	// Initialize Storage
	store, err := storage.New(storage.Options{
		StorageType: cfg.App.Storage,
		Logg:        logg,
	})
	if err != nil {
		log.Fatalf("failed to init storage: %v", err)
	}

	if err := store.Connect(); err != nil {
		logg.Error("failed to connect storage: ", "error", err.Error())
		os.Exit(1)
	}

	// Initialize Repository
	calendar := app.New(logg, store)

	// Initialize Server
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

	// Create context for signal handling
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	logg.Info("initializing calendar server...",
		slog.String("server_type", cfg.App.Server),
		slog.String("storage_type", cfg.App.Storage),
	)

	// Start server with the context that handles signals
	if err := srv.Start(ctx); err != nil {
		logg.Error("failed to start server",
			slog.String("error", err.Error()),
			slog.String("server_type", cfg.App.Server),
		)
		os.Exit(1)
	}

	// wait for context cancellation
	<-ctx.Done()
	logg.Info("received shutdown signal")

	// create new context for shutdown with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Second*3)
	defer shutdownCancel()

	// start graceful shutdown
	shutdownErr := make(chan error, 1)
	go func() {
		// First stop the server
		if err := srv.Stop(shutdownCtx); err != nil {
			shutdownErr <- err
			return
		}
		// Then close the storage
		if err := store.Close(); err != nil {
			shutdownErr <- err
			return
		}
		shutdownErr <- nil
	}()

	// wait for shutdown to complete or timeout
	select {
	case err := <-shutdownErr:
		if err != nil {
			logg.Error("graceful shutdown failed",
				slog.String("error", err.Error()),
				slog.String("server_type", cfg.App.Server),
			)
			os.Exit(1)
		}

		logg.Info("server stopped gracefully",
			slog.String("server_type", cfg.App.Server),
			slog.String("status", "success"),
		)
		os.Exit(0)

	case <-shutdownCtx.Done():
		logg.Error("graceful shutdown timed out",
			slog.String("timeout", "3s"),
			slog.String("action", "forcing exit"),
		)
		os.Exit(1)
	}
}
