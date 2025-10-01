package oak

import (
	"context"
	"log"
	"log/slog"
	"os"
)

func init() {
	logger := slog.New(TintHandler(os.Stdout, LevelTrace))
	slog.SetDefault(logger)
}

func SetDefault(l *Oak) {
	slog.SetDefault(l.logger)
}

func NewLogLogger(l *Oak, level slog.Level) *log.Logger {
	return slog.NewLogLogger(l.logger.Handler(), level)
}

func Info(msg string, args ...any) {
	slog.Info(msg, args...)
}

func Error(msg string, args ...any) {
	slog.Error(msg, args...)
}

func WithContext(ctx context.Context, l *Oak) context.Context {
	return context.WithValue(ctx, CtxName, l)
}

func FromContext(ctx context.Context, layer string) *Oak {
	v, ok := ctx.Value(CtxName).(*Oak)
	if !ok {
		logger := New(TintHandler(os.Stdout, LevelTrace)).Layer(layer)
		logger.Error("no logger found in context, new one created")
		return logger
	}
	return v.Layer(layer)
}
