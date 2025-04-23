package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/internal/app"
	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/pkg/logger"
	"github.com/streadway/amqp"
)

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

	// Подписка на очередь
	msgs, err := channel.Consume(
		cfg.RabbitMQ.Queue, // имя очереди
		"",                 // consumer
		true,               // auto-ack
		false,              // exclusive
		false,              // no-local
		false,              // no-wait
		nil,                // args
	)
	if err != nil {
		logg.Error("failed to consume messages", "error", err)
		return
	}

	logg.Info("RabbitMQ consumer is ready", "queue", cfg.RabbitMQ.Queue)

	// Обработка сигналов для безопасной остановки
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	done := make(chan bool)

	go func() {
		for msg := range msgs {
			var event app.Event
			if err := json.Unmarshal(msg.Body, &event); err != nil {
				logg.Warn("failed to unmarshal message", "error", err)
				continue
			}

			// Логируем уведомление
			logg.Info(
				"notification received",
				"event_id", event.ID,
				"event_title", event.Title,
				"event_time", event.StartTime,
			)
		}
	}()

	go func() {
		<-signalChan
		logg.Info("Shutting down sender")
		done <- true
	}()

	<-done
}
