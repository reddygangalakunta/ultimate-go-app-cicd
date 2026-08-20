package logger

import (
	"log/slog"
	"os"
	"strings"
)

// SetupLogger initializes a structured JSON logger using log/slog.
func SetupLogger(level string) *slog.Logger {
	var logLevel slog.Level

	switch strings.ToLower(level) {
	case "debug":
		logLevel = slog.LevelDebug
	case "warn", "warning":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: logLevel,
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	l := slog.New(handler)
	slog.SetDefault(l)

	return l
}
