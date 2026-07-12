package logger

import (
	"io"
	"log/slog"
)

type Format string

const (
	Console Format = "console"
	Text    Format = "text"
	JSON    Format = "json"
)

type Config struct {
	Level  slog.Level
	Format Format
	Output io.Writer
}

func New(cfg Config) *slog.Logger {
	handler := newHandler(cfg)

	return slog.New(handler)
}
