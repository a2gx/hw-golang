package internalhttp

import (
	"fmt"
	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/pkg/logger"
	"net/http"
	"time"
)

type MiddlewareFn func(http.Handler) http.Handler

func applyMiddleware(h http.Handler, mws ...MiddlewareFn) http.Handler {
	for _, mw := range mws {
		h = mw(h)
	}
	return h
}

func LoggingMiddleware(logg logger.Logger) MiddlewareFn {
	return func(next http.Handler) http.Handler {
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

			logg.Info(msg)
		})
	}
}
