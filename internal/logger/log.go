package logger

import (
	"context"
	"log/slog"
	"os"
)

const (
	OperationField = "operation"
	ErrFiled       = "err"
)

type Wrapper struct {
	logger *slog.Logger
}

// New returns a JSON logger that writes debug-level logs to standard output.
func New() *Wrapper {
	l := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:       slog.LevelDebug,
		ReplaceAttr: replaceTimeFormat,
	}))
	return &Wrapper{
		logger: l,
	}
}

// replaceTimeFormat shortens log timestamps while keeping their date and minute.
func replaceTimeFormat(_ []string, attr slog.Attr) slog.Attr {
	if attr.Key == slog.TimeKey {
		attr.Value = slog.StringValue(attr.Value.Time().Format("2006-01-02|15:04"))
	}
	return attr
}

// With returns a logger with the given attribute attached.
func (w *Wrapper) With(key, value string) *Wrapper {
	return &Wrapper{
		logger: w.logger.With(key, value),
	}
}

// Debug logs a message at debug level.
func (w *Wrapper) Debug(ctx context.Context, msg string, args ...any) {
	w.logger.DebugContext(ctx, msg, args...)
}

// Info logs a message at info level.
func (w *Wrapper) Info(ctx context.Context, msg string, args ...any) {
	w.logger.InfoContext(ctx, msg, args...)
}

// Warn logs a message at warning level.
func (w *Wrapper) Warn(ctx context.Context, msg string, args ...any) {
	w.logger.WarnContext(ctx, msg, args...)
}

// Error logs a message at error level.
func (w *Wrapper) Error(ctx context.Context, msg string, args ...any) {
	w.logger.ErrorContext(ctx, msg, args...)
}
