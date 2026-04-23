package logging

import (
	"context"
	"log/slog"
)

type contextLoggerKey struct{}

// ContextWithLogger returns ctx with logger attached.
func ContextWithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, contextLoggerKey{}, logger)
}

// FromContext returns the logger attached to ctx, or slog.Default() if none.
func FromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(contextLoggerKey{}).(*slog.Logger); ok {
		return logger
	}
	return slog.Default()
}
