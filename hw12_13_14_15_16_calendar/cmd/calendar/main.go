package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/internal/app"
	internalhttp "github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/internal/server/http"
	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/internal/storage"
	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/pkg/logger"
)

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	// initial Config
	config, err := NewConfig()
	if err != nil {
		log.Fatal("filed to load config:", err)
	}

	// initial Logger
	logg := logger.New(logger.Options{
		Level:    config.Logger.Level,
		Handler:  config.Logger.Handler,
		Filename: config.Logger.Filename,
		Source:   config.Logger.Source,
	})
	defer logg.Close()

	// initial Storage
	store, err := storage.New(storage.Options{
		StorageType: config.Storage.Source,
	})
	if err != nil {
		log.Fatalf("failed to init storage: %v", err)
	}

	// initial Business-logic
	calendar := app.New(logg, store)

	// initial Server
	server := internalhttp.NewServer(logg, calendar)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
