package serverhttp

import (
	"fmt"
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

		msg := fmt.Sprintf(
			"%s [%s] %s %s %s %dms '%s'",
			r.RemoteAddr,
			start.Format("02/Jan/2006:15:04:05 -0700"),
			r.Method,
			r.RequestURI,
			r.Proto,
			latency.Milliseconds(),
			r.UserAgent(),
		)

		slog.Info(msg)
	})
}
