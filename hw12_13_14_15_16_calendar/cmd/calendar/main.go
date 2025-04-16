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

	// Инициализация конфигурации
	cfg, err := NewConfig()
	if err != nil {
		log.Fatalf("failed to init configuration: %v", err)
	}

	// Инициализация логгера
	logg := logger.New(logger.Options{
		Level:     cfg.Logger.Level,
		Handler:   cfg.Logger.Handler,
		Filename:  cfg.Logger.Filename,
		AddSource: cfg.Logger.AddSource,
	})

	// Инициализация хранилища
	store, err := storage.New(storage.Options{
		StorageType: cfg.App.Storage,
		Logg:        logg,
	})
	if err != nil {
		logg.Error("failed to init storage", "error", err)
		os.Exit(1)
	}
	if err := store.Connect(); err != nil {
		logg.Error("failed to connect storage", "error", err)
		os.Exit(1)
	}

	// Инициализация приложения
	calendar := app.New(logg, store)

	// Инициализация сервера
	srv, err := server.New(server.Options{
		ServerType: cfg.App.Server,
		HttpAddr:   cfg.Server.HttpAddr,
		GrpcAddr:   cfg.Server.GrpcAddr,
		Logg:       logg,
		App:        calendar,
	})
	if err != nil {
		logg.Error("failed to init server", "error", err)
		os.Exit(1)
	}

	// Создание контекста для обработки сигналов
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	// Каналы для обработки завершения работы
	shutdownErr := make(chan error, 1)
	shutdownDone := make(chan struct{})

	// Обработка shutdown сигналов
	go func() {
		<-ctx.Done()
		logg.Warn("received shutdown signal")

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := srv.Stop(ctx); err != nil {
			shutdownErr <- err
			return
		}

		if err := store.Close(); err != nil {
			shutdownErr <- err
			return
		}

		close(shutdownDone)
	}()

	logg.Info("calendar is running...",
		"server_type", cfg.App.Server,
		"storage_type", cfg.App.Storage,
	)

	// Запуск сервера
	if err := srv.Start(ctx); err != nil {
		shutdownErr <- err
	}

	// Ожидание завершения работы
	select {
	case err := <-shutdownErr:
		logg.Error("shutdown error", "error", err)
		cancel()
		os.Exit(1)
	case <-shutdownDone:
		logg.Info("application shutdown completed successfully")
		os.Exit(0)
	}
}
