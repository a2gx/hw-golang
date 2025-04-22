package serverhttp

import (
	"log/slog"
	"net/http"
	"time"
)

type nextFn func(http.Handler) http.Handler

func applyMiddleware(next http.Handler, appNextFn ...nextFn) http.Handler {
	for _, nFn := range appNextFn {
		next = nFn(next)
	}
	return next
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		latency := time.Since(start)

		slog.Info(
			"HTTP request",
			"remote_addr", r.RemoteAddr,
			"time", start.Format("02/Jan/2006:15:04:05 -0700"),
			"method", r.Method,
			"request_uri", r.RequestURI,
			"proto", r.Proto,
			"latency_ms", latency.Milliseconds(),
			"user_agent", r.UserAgent(),
		)
	})
}
