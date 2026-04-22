// Package logging provides application-wide slog configuration.
//
// It is self-contained and intended to be copied into each dnsimple Go
// application's internal/ directory. No project-specific imports.
package logging

import (
	"log/slog"
	"os"
	"strings"
)

// New returns a configured slog.Logger that writes JSON to stdout at the
// given level. It does not mutate the global slog state — callers decide
// whether to call slog.SetDefault on the result.
func New(level slog.Level) *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
	}))
}

// Err returns an slog.Attr with the canonical "error" key.
//
// Use this helper for all error logging so the field name is consistent across
// the codebase and callers can't mistakenly use a different key ("err", "errno",
// "cause", etc.).
//
//	slog.Error("Failed to deliver", logging.Err(err))
func Err(err error) slog.Attr {
	return slog.Any("error", err)
}

// ParseLevel maps a string to an slog.Level.
// Recognized values (case-insensitive): "debug", "info", "warn", "error".
// Any other value (including empty) returns slog.LevelInfo.
func ParseLevel(s string) slog.Level {
	switch strings.ToLower(s) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
