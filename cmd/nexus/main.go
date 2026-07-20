package main

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/shubhindia/nexus/internal/config"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/shubhindia/nexus/internal/api"
	"github.com/shubhindia/nexus/internal/gateway"
	"github.com/shubhindia/nexus/internal/logger"
	"github.com/shubhindia/nexus/internal/metrics"
	"github.com/shubhindia/nexus/internal/provider"
	"github.com/shubhindia/nexus/internal/providers/gemini"
	"github.com/shubhindia/nexus/internal/providers/llamacpp"
)

func main() {
	cfg := config.Load()

	// Logger
	log := logger.New(logger.Config{
		Format: logger.Console,
		Level:  slog.LevelInfo,
	})

	registry := prometheus.NewRegistry()

	// Metrics
	m := metrics.New(registry)

	// LLM Provider
	p, err := newProvider(cfg)
	if err != nil {
		log.Error(
			"provider.init",
			slog.Any("error", err),
		)
		return
	}

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
		slog.String("addr", cfg.NexusPort),
		slog.String("provider", p.Name()),
	)

	if err := http.ListenAndServe(cfg.NexusPort, mux); err != nil {
		log.Error(
			"http.server.error",
			slog.Any("error", err),
		)
	}
}

func newProvider(cfg config.Config) (provider.Provider, error) {
	switch cfg.Provider {
	case config.ProviderLlamaCPP:
		return llamacpp.New(cfg.LlamaURL), nil
	case config.ProviderGemini:
		if cfg.GeminiKey == "" {
			return nil, fmt.Errorf("NEXUS_GEMINI_API_KEY is required for provider %q", cfg.Provider)
		}

		return gemini.New(cfg.GeminiURL, cfg.GeminiKey), nil
	default:
		return nil, fmt.Errorf("unsupported provider %q", cfg.Provider)
	}
}
