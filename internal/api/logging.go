package api

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/shubhindia/nexus/internal/logger"
)

func Logging(base *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestLogger := base.With(
				slog.String("request_id", RequestIDFromContext(r.Context())),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
			)

			ctx := logger.WithContext(r.Context(), requestLogger)

			start := time.Now()

			requestLogger.Info("http.request.started")

			rw := newResponseWriter(w)

			next.ServeHTTP(rw, r.WithContext(ctx))

			requestLogger.Info(
				"http.request.completed",
				slog.Int("status", rw.status),
				slog.Int("bytes", rw.bytes),
				slog.Duration("duration", time.Since(start)),
			)
		})
	}
}
