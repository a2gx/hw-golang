package logger

import (
	"io"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/lmittmann/tint"
)

type Logger struct {
	slog.Logger
	closeFn
}

type Options struct {
	Level     string
	Handler   string
	Filename  string
	AddSource bool
}

type internalOptions struct {
	Level slog.Level
	Options
}

type closeFn func()

func New(opts Options) *Logger {
	writer, closeFn := createWriter(opts)
	handler := createHandler(writer, internalOptions{
		Level:   parseLevel(opts.Level),
		Options: opts,
	})

	logger := slog.New(handler)
	slog.SetDefault(logger) // set logger globally

	return &Logger{
		Logger:  *logger,
		closeFn: closeFn,
	}
}

func (l *Logger) Close() {
	if l.closeFn != nil {
		l.closeFn()
	}
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
		log.Printf("Unknown log level '%s', defaulting to 'info'", lvl)
		return slog.LevelInfo
	}
}

func createWriter(opts Options) (io.Writer, closeFn) {
	var defaultCloseFn = func() {}

	if opts.Filename == "" {
		return os.Stdout, defaultCloseFn
	}

	writer, err := os.OpenFile(opts.Filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		log.Printf("Failed to open log file: %v", err)
		return os.Stdout, defaultCloseFn
	}

	return writer, func() {
		if err := writer.Close(); err != nil {
			log.Printf("Failed to close log file: %v", err)
		}
	}
}

func createHandler(w io.Writer, opts internalOptions) slog.Handler {
	if opts.Handler == "" {
		opts.Handler = "text" // set default handler
	}

	switch strings.ToLower(opts.Handler) {
	case "text_color":
		return tint.NewHandler(w, &tint.Options{
			Level:     opts.Level,
			AddSource: opts.AddSource,
		})
	case "text":
		return slog.NewTextHandler(w, &slog.HandlerOptions{
			Level:     opts.Level,
			AddSource: opts.AddSource,
		})
	case "json":
		return slog.NewJSONHandler(w, &slog.HandlerOptions{
			Level:     opts.Level,
			AddSource: opts.AddSource,
		})
	default:
		log.Printf("Unknown handler '%s', defaulting to 'text'", opts.Handler)
		return slog.NewTextHandler(w, &slog.HandlerOptions{
			Level:     opts.Level,
			AddSource: opts.AddSource,
		})
	}
}
