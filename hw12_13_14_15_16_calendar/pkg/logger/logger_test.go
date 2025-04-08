package logger

import (
	"bytes"
	"log/slog"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	log := New(Options{Level: "debug"})
	if log == nil {
		t.Error("Logger is nil")
		return
	}

	log.Close()
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected slog.Level
	}{
		{"debug", slog.LevelDebug},
		{"info", slog.LevelInfo},
		{"warn", slog.LevelWarn},
		{"error", slog.LevelError},
	}

	for _, test := range tests {
		if got := parseLevel(test.input); got != test.expected {
			t.Errorf("parseLevel(%q) = %v; want %v", test.input, got, test.expected)
		}
	}
}

func TestCreateWriter_File(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "log_test_*.log")
	if err != nil {
		t.Errorf("Failed to create temporary file: %v", err)
	}
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	opts := Options{Filename: tmpFile.Name()}
	writer, closeFn := createWriter(opts)

	if writer == nil {
		t.Error("Writer is nil")
	}

	closeFn()
}

func TestCreateWriter_Stdout(t *testing.T) {
	opts := Options{}
	writer, closeFn := createWriter(opts)

	if writer != os.Stdout {
		t.Error("Writer is not stdout")
	}
	if closeFn == nil {
		t.Error("Close function is nil")
	} else {
		closeFn()
	}
}

func TestCreateHandler(t *testing.T) {
	var buf bytes.Buffer

	tests := []struct {
		handlerType string
	}{
		{"text"},
		{"json"},
		{"text_color"},
	}

	for _, tt := range tests {
		t.Run(tt.handlerType, func(t *testing.T) {
			opts := internalOptions{
				Level: slog.LevelInfo,
				Options: Options{
					HandlerType: tt.handlerType,
				},
			}
			handler := createHandler(&buf, opts)
			if handler == nil {
				t.Fatalf("handler should not be nil: %v", handler)
			}
		})
	}
}
