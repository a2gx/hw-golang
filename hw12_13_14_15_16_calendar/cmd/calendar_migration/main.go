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
	flag.Parse()

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

	// Открываем соединение с базой данных
	db, err := goose.OpenDBWithDriver("postgres", cfg.DatabaseDns)
	if err != nil {
		logg.Error("Не удалось подключиться к базе данных: %v", err)
		os.Exit(1)
	}
	defer db.Close()

	switch cfg.MigrationCommand {
	case "up":
		if err := goose.Up(db, cfg.MigrationDir); err != nil {
			logg.Error("Ошибка при применении миграций: %v", err)
			os.Exit(1)
		}
	case "down":
		if err := goose.Down(db, cfg.MigrationDir); err != nil {
			logg.Error("Ошибка при откате миграций: %v", err)
			os.Exit(1)
		}
	case "status":
		if err := goose.Status(db, cfg.MigrationDir); err != nil {
			logg.Error("Ошибка при получении статуса миграций: %v", err)
			os.Exit(1)
		}
	case "create":
		if cfg.MigrationName == "" {
			logg.Error("Для команды create необходимо указать имя миграции с помощью флага -name")
			os.Exit(1)
		}
		if err := goose.Create(nil, cfg.MigrationDir, cfg.MigrationName, "sql"); err != nil {
			logg.Error("Ошибка при создании новой миграции: %v", err)
			os.Exit(1)
		}
	default:
		logg.Error("Неизвестная команда: %s", cfg.MigrationCommand)
	}

	logg.Info("Миграции успешно выполнены")
}
