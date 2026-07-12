package logger

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

func newHandler(cfg Config) slog.Handler {
	output := cfg.Output
	if output == nil {
		output = os.Stderr
	}

	opts := &slog.HandlerOptions{
		Level: cfg.Level,
	}

	switch cfg.Format {
	case JSON:
		return slog.NewJSONHandler(output, opts)

	case Text:
		return slog.NewTextHandler(output, opts)

	case Console:
		fallthrough

	default:
		return tint.NewHandler(output, &tint.Options{
			Level: cfg.Level,
		})
	}
}
