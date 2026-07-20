package main

import (
	"log/slog"
	"net/http"

	"github.com/shubhindia/nexus/internal/config"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/shubhindia/nexus/internal/api"
	"github.com/shubhindia/nexus/internal/gateway"
	"github.com/shubhindia/nexus/internal/logger"
	"github.com/shubhindia/nexus/internal/metrics"
	"github.com/shubhindia/nexus/internal/providers/llamacpp"
)

func main() {
	// Logger
	log := logger.New(logger.Config{
		Format: logger.Console,
		Level:  slog.LevelInfo,
	})

	registry := prometheus.NewRegistry()

	// Metrics
	m := metrics.New(registry)

	// LLM Provider
	p := llamacpp.New(config.LlamaURL)

	gw := gateway.New(p)

	// API Router
	router := api.NewRouter(
		gw,
		log,
		m,
	)

	// Root mux
	mux := http.NewServeMux()

	mux.Handle("/", router)
	mux.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	log.Info(
		"http.server.start",
		slog.String("addr", config.NexusPort),
	)

	if err := http.ListenAndServe(config.NexusPort, mux); err != nil {
		log.Error(
			"http.server.error",
			slog.Any("error", err),
		)
	}
}
