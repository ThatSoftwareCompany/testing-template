package logging

import (
	"log/slog"
	"os"
	"strings"
)

func New(level, environment string) *slog.Logger {
	configuredLevel := slog.LevelInfo
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug":
		configuredLevel = slog.LevelDebug
	case "warn":
		configuredLevel = slog.LevelWarn
	case "error":
		configuredLevel = slog.LevelError
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: configuredLevel,
	})).With("environment", environment)
}
