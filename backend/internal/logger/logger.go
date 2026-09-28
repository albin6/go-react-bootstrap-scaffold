package logger

import (
	"context"
	"log/slog"
	"os"
	"strings"
)

type contextKey struct{}

var loggerContextKey = contextKey{}

// Init initializes the default slog logger based on level and format.
func Init(levelStr, format string) *slog.Logger {
	var level slog.Level
	switch strings.ToLower(levelStr) {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	var handler slog.Handler
	if strings.ToLower(format) == "json" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	l := slog.New(handler)
	slog.SetDefault(l)
	return l
}

// WithContext returns a new context with the logger attached.
func WithContext(ctx context.Context, l *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerContextKey, l)
}

// FromContext extracts the logger from context or returns the default slog logger.
func FromContext(ctx context.Context) *slog.Logger {
	if ctx != nil {
		if l, ok := ctx.Value(loggerContextKey).(*slog.Logger); ok && l != nil {
			return l
		}
	}
	return slog.Default()
}
