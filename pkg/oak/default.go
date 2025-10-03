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

func Trace(msg string, args ...any) {
	slog.Log(context.TODO(), LevelTrace, msg, args...)
}

func Debug(msg string, args ...any) {
	slog.Debug(msg, args...)
}

func Info(msg string, args ...any) {
	slog.Info(msg, args...)
}

func Notice(msg string, args ...any) {
	slog.Log(context.TODO(), LevelNotice, msg, args...)
}

func Warning(msg string, args ...any) {
	slog.Log(context.TODO(), LevelWarning, msg, args...)
}

func Error(msg string, args ...any) {
	slog.Error(msg, args...)
}

func Fatal(msg string, args ...any) {
	slog.Log(context.TODO(), LevelFatal, msg, args...)
}

func WithContext(ctx context.Context, l *Oak) context.Context {
	return context.WithValue(ctx, CtxName, l)
}

func FromContextWithLayer(ctx context.Context, layer string) *Oak {
	return fromContextWithLayer(ctx, layer)
}

func FromContext(ctx context.Context) *Oak {
	return fromContextWithLayer(ctx, "")
}

func fromContextWithLayer(ctx context.Context, layer string) *Oak {
	v, ok := ctx.Value(CtxName).(*Oak)
	if !ok {
		logger := New(TintHandler(os.Stdout, LevelTrace)).Layer(layer)
		logger.Error("no logger found in context, new one created")
		return logger
	}
	return v.Layer(layer)
}
