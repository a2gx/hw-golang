package logger

// полезности:
// https://betterstack.com/community/guides/logging/logging-in-go/

import (
	"io"
	"log/slog"
	"strings"
)

func New(level string, writer io.Writer) *slog.Logger {
	//Обычный вывод JSON в файл
	handler := slog.NewJSONHandler(writer, &slog.HandlerOptions{
		AddSource: true,
		Level:     parseLevel(level),
	})

	//Красивый вывод в КОНСОЛЬ!
	//handler := tint.NewHandler(os.Stdout, &tint.Options{
	//	AddSource: true,
	//	Level:     parseLevel(level),
	//})

	logger := slog.New(handler)
	slog.SetDefault(logger) // глобально доступная конфигурация !!!

	return logger
}

func parseLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
