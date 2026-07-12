package api

import (
	"log/slog"
	"net/http"

	"github.com/shubhindia/nexus/internal/gateway"
	"github.com/shubhindia/nexus/internal/metrics"
)

func NewRouter(
	gw *gateway.Gateway,
	log *slog.Logger,
	m *metrics.Metrics,
) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", Health)
	mux.HandleFunc("/v1/models", Models(gw))
	mux.HandleFunc(
		"/v1/chat/completions",
		Chat(gw),
	)

	return Chain(
		mux,
		RequestID,
		Logging(log),
		metrics.Middleware(m),
	)
}
