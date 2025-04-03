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
	closeFn closeFn
	//Close   func()
}

type Options struct {
	Handler  string
	Filename string
	Source   bool
}

type internalOptions struct {
	Level slog.Level
	Options
}

type closeFn func()

func New(level string, opts ...Options) *Logger {
	options := parseOptions(opts)
	writer, closeFn := createWriter(options)
	handler := createHandler(writer, internalOptions{
		Level:   parseLevel(level),
		Options: options,
	})

	logger := slog.New(handler)
	slog.SetDefault(logger)

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

func parseOptions(opts []Options) Options {
	if len(opts) == 1 {
		return opts[0]
	}
	return Options{}
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
		log.Printf("Ooops... Unknown log level '%s', defaulting to 'info'", lvl)
		return slog.LevelInfo
	}
}

func createWriter(o Options) (io.Writer, closeFn) {
	var defaultCloseFn closeFn = func() {}

	if len(o.Filename) == 0 {
		return os.Stdout, defaultCloseFn
	}

	writer, err := os.OpenFile(o.Filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Ooops... Failed to open log file: %v", err)
		return os.Stdout, defaultCloseFn
	}

	return writer, func() {
		if err := writer.Close(); err != nil {
			log.Printf("Ooops... Failed to close log file: %v", err)
		}
	}
}

func createHandler(w io.Writer, o internalOptions) slog.Handler {
	if o.Handler == "" {
		o.Handler = "text"
	}

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
