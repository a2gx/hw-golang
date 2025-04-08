package main

import (
	"context"
	"flag"
	"fmt"
	"log"
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
	serv, err := server.New(server.Options{
		ServerType: cfg.App.Server,

		AddrHttp: fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		AddrGrpc: fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),

		Logg: logg,
		App:  calendar,
	})
	if err != nil {
		log.Fatalf("failed to init server: %v", err)
	}

	// Start server ...

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		_, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		if err := serv.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	if err := serv.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1)
	}
}
