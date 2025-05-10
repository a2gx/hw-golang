package main

import (
	"flag"
	"log"
	"os"

	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/pkg/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func main() {
	exitCode := 0
	defer func() {
		os.Exit(exitCode)
	}()

	flag.Parse()

	// Инициализация конфигурации
	cfg, err := NewConfig()
	if err != nil {
		log.Printf("failed to init configuration: %v", err)
		return
	}

	// Инициализация логгера
	logg := logger.New(logger.Options{
		Level:     cfg.Logger.Level,
		Handler:   cfg.Logger.Handler,
		Filename:  cfg.Logger.Filename,
		AddSource: cfg.Logger.AddSource,
	})

	// Открываем соединение с базой данных
	db, err := goose.OpenDBWithDriver("postgres", cfg.DatabaseDNS)
	if err != nil {
		logg.Error("failed to connect to database: " + err.Error())
		exitCode = 1
		return
	}
	defer db.Close()

	switch cfg.MigrationCommand {
	case "up":
		if err := goose.Up(db, cfg.MigrationDir); err != nil {
			logg.Error("failed to apply migrations: " + err.Error())
			exitCode = 1
			return
		}
	case "down":
		if err := goose.Down(db, cfg.MigrationDir); err != nil {
			logg.Error("failed to rollback migrations: " + err.Error())
			exitCode = 1
			return
		}
	case "status":
		if err := goose.Status(db, cfg.MigrationDir); err != nil {
			logg.Error("failed to get migration status: " + err.Error())
			exitCode = 1
			return
		}
	case "create":
		if cfg.MigrationName == "" {
			logg.Error("migration name is required for create command (use -name flag)")
			exitCode = 1
			return
		}
		if err := goose.Create(nil, cfg.MigrationDir, cfg.MigrationName, "sql"); err != nil {
			logg.Error("failed to create new migration: " + err.Error())
			exitCode = 1
			return
		}
	default:
		logg.Error("unknown command: " + cfg.MigrationCommand)
		exitCode = 1
		return
	}
}
