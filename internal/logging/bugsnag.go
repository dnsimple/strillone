package logging

import (
	"fmt"
	"log/slog"
)

// BugsnagLogger adapts an *slog.Logger to the Printf-based Logger interface
// expected by github.com/bugsnag/bugsnag-go.
//
// Bugsnag uses Printf only for its own internal diagnostics (event delivery
// failures, configuration warnings), not for customer error reports. We emit
// those at Warn because they typically indicate a problem with error
// reporting itself.
type BugsnagLogger struct {
	Logger *slog.Logger
}

// Printf implements the Bugsnag Logger interface.
func (b *BugsnagLogger) Printf(format string, args ...any) {
	b.Logger.Warn(fmt.Sprintf(format, args...))
}
