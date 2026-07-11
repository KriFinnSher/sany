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

func New() *Wrapper {
	l := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	return &Wrapper{
		logger: l,
	}
}

func (w *Wrapper) With(key, value string) *Wrapper {
	return &Wrapper{
		logger: w.logger.With(key, value),
	}
}

func (w *Wrapper) Debug(ctx context.Context, msg string, args ...any) {
	w.logger.DebugContext(ctx, msg, args...)
}

func (w *Wrapper) Info(ctx context.Context, msg string, args ...any) {
	w.logger.InfoContext(ctx, msg, args...)
}

func (w *Wrapper) Warn(ctx context.Context, msg string, args ...any) {
	w.logger.WarnContext(ctx, msg, args...)
}

func (w *Wrapper) Error(ctx context.Context, msg string, args ...any) {
	w.logger.ErrorContext(ctx, msg, args...)
}
