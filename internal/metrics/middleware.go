package metrics

import (
	"net/http"
	"strconv"
	"time"
)

type responseWriter struct {
	http.ResponseWriter

	status int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		status:         http.StatusOK,
	}
}

func (w *responseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func Middleware(m *Metrics) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rw := newResponseWriter(w)

			start := time.Now()

			next.ServeHTTP(rw, r)

			status := strconv.Itoa(rw.status)

			m.HTTP.Requests.
				WithLabelValues(
					r.Method,
					r.URL.Path,
					status,
				).
				Inc()

			m.HTTP.Duration.
				WithLabelValues(
					r.Method,
					r.URL.Path,
				).
				Observe(time.Since(start).Seconds())
		})
	}
}
