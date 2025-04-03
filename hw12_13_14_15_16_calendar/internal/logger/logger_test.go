package logger

import (
	"log/slog"
	"os"
	"testing"
)

func TestLogger(t *testing.T) {
	logg, closeFn := New("debug", Options{
		Handler: "json",
	})
	defer closeFn()

	if logg == nil {
		t.Fatal("logger is nil")
	}
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

func TestCreateWriter(t *testing.T) {
	w, _ := createWriter(Options{Handler: "text"})
	if w != os.Stdout {
		t.Errorf("createWriter({handler}) = %v; want os.Stdout", w)
	}
}
