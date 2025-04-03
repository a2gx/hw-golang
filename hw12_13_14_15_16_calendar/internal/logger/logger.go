package logger

import (
	"io"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/lmittmann/tint"
)

type Options struct {
	Handler  string
	Filename string
	Source   bool
}

type internalOptions struct {
	Level slog.Level
	Options
}

type WriterClose func()

func New(level string, opts ...Options) (*slog.Logger, WriterClose) {
	options := parseOptions(opts) // а говорили в Go все аргументы обязательны XD
	writer, c := createWriter(options)
	handler := createHandler(writer, internalOptions{
		Level:   parseLevel(level),
		Options: options,
	})

	logger := slog.New(handler)
	slog.SetDefault(logger)

	return logger, c
}

func parseLevel(lvl string) slog.Level {
	switch strings.ToLower(lvl) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		log.Printf("Ooops... Uncnown log level '%s', defaulting to 'info'", lvl)
		return slog.LevelInfo
	}
}

func parseOptions(opts []Options) Options {
	if len(opts) == 1 {
		return opts[0]
	}
	return Options{}
}

func createWriter(opts Options) (io.Writer, WriterClose) {
	var defaultClose WriterClose = func() {}

	if len(opts.Filename) == 0 {
		return os.Stdout, defaultClose
	}

	writer, err := os.OpenFile(opts.Filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Ooops... Failed to open log file: %v", err)
		return os.Stdout, defaultClose
	}

	return writer, func() {
		if err := writer.Close(); err != nil {
			log.Printf("Ooops... Failed to close log file: %v", err)
		}
	}
}

func createHandler(w io.Writer, o internalOptions) slog.Handler {
	switch strings.ToLower(o.Handler) {
	case "text_color":
		return tint.NewHandler(w, &tint.Options{
			Level:     o.Level,
			AddSource: o.Source,
		})
	case "text":
		return slog.NewTextHandler(w, &slog.HandlerOptions{
			Level:     o.Level,
			AddSource: o.Source,
		})
	case "json":
		return slog.NewJSONHandler(w, &slog.HandlerOptions{
			Level:     o.Level,
			AddSource: o.Source,
		})
	default:
		log.Printf("Ooops... Unknown handler '%s', defaulting to 'text'", o.Handler)
		return slog.NewTextHandler(w, &slog.HandlerOptions{
			Level:     o.Level,
			AddSource: o.Source,
		})
	}
}
