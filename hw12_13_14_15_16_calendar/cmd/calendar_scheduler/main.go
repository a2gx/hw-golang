package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/internal/app"
	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/internal/storage"
	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/pkg/logger"
	"github.com/streadway/amqp"
)

func sendNotification(channel *amqp.Channel, queueName string, event app.Event) error {
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	err = channel.Publish(
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

func main() {
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

	// Инициализация хранилища
	store, err := storage.New(storage.Options{
		StorageType: cfg.App.Storage,
		DatabaseDNS: cfg.DatabaseDNS,
		Logg:        logg,
	})
	if err != nil {
		logg.Error("failed to init storage", "error", err)
		return
	}
	if err := store.Connect(); err != nil {
		logg.Error("failed to connect storage", "error", err)
		return
	}
	defer store.Close()

	// Подключение к RabbitMQ
	rabbitConn, err := amqp.Dial(cfg.RabbitDNS)
	if err != nil {
		logg.Error("failed to connect to RabbitMQ", "error", err)
		return
	}
	defer rabbitConn.Close()

	channel, err := rabbitConn.Channel()
	if err != nil {
		logg.Error("failed to open a channel", "error", err)
		return
	}
	defer channel.Close()

	// Создание очереди
	queue, err := channel.QueueDeclare(
		cfg.RabbitMQ.Queue, // имя очереди
		true,               // durable
		false,              // auto-delete
		false,              // exclusive
		false,              // no-wait
		nil,                // arguments
	)
	if err != nil {
		logg.Error("failed to declare a queue", "error", err)
		return
	}

	logg.Info("RabbitMQ queue is ready", "queue", cfg.RabbitMQ.Queue)

	// Обработка сигналов для безопасной остановки
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signalChan
		cancel()
	}()

	// Запуск планировщика
	ticker := time.NewTicker(time.Duration(cfg.App.TimeoutScheduler) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logg.Info("Shutting down scheduler")
			return
		case <-ticker.C:
			logg.Info("Scheduler tick")

			// Выбор событий для уведомления
			events, err := store.FetchEventsToNotify()
			if err != nil {
				logg.Warn("failed to fetch events", "error", err)
				continue
			}

			// Отправка уведомлений в RabbitMQ
			for _, event := range events {
				if err := sendNotification(channel, queue.Name, event); err != nil {
					logg.Warn("failed to send notification", "error", err)
				} else {
					logg.Debug(
						"notification sent",
						"event_id", event.ID,
						"event_title", event.Title,
					)
				}
			}

			// Удаление старых событий
			if err := store.DeleteOldEvents(); err != nil {
				logg.Warn("failed to delete old events", "error", err)
			}
		}
	}
}
