package server_http

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/internal/app"
	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/pkg/logger"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	logg := logger.New(logger.Options{Level: "debug"})
	application := &app.App{}

	server := New(":8080", logg, application)
	require.NotNil(t, server)
	require.Equal(t, ":8080", server.srv.Addr)
}

func TestStartAndStop(t *testing.T) {
	logg := logger.New(logger.Options{Level: "debug"})
	application := &app.App{}

	server := New(":8081", logg, application)
	require.NotNil(t, server)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		err := server.Start(ctx)
		require.NoError(t, err)
	}()

	// Проверка доступности эндпоинта /ping
	time.Sleep(100 * time.Millisecond) // Даем серверу время на запуск
	resp, err := http.Get("http://localhost:8081/ping")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Остановка сервера
	err = server.Stop(ctx)
	require.NoError(t, err)
}
