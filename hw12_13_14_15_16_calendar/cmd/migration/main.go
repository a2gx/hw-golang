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
		logg.Error("Не удалось подключиться к базе данных: " + err.Error())
		exitCode = 1
		return
	}
	defer db.Close()

	switch cfg.MigrationCommand {
	case "up":
		if err := goose.Up(db, cfg.MigrationDir); err != nil {
			logg.Error("Ошибка при применении миграций: " + err.Error())
			exitCode = 1
			return
		}
	case "down":
		if err := goose.Down(db, cfg.MigrationDir); err != nil {
			logg.Error("Ошибка при откате миграций: " + err.Error())
			exitCode = 1
			return
		}
	case "status":
		if err := goose.Status(db, cfg.MigrationDir); err != nil {
			logg.Error("Ошибка при получении статуса миграций: " + err.Error())
			exitCode = 1
			return
		}
	case "create":
		if cfg.MigrationName == "" {
			logg.Error("Для команды create необходимо указать имя миграции с помощью флага -name")
			exitCode = 1
			return
		}
		if err := goose.Create(nil, cfg.MigrationDir, cfg.MigrationName, "sql"); err != nil {
			logg.Error("Ошибка при создании новой миграции: " + err.Error())
			exitCode = 1
			return
		}
	default:
		logg.Error("Неизвестная команда: " + cfg.MigrationCommand)
		exitCode = 1
		return
	}

	logg.Info("Миграции успешно выполнены")
}
