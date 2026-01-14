package oak

import (
	"context"
	"log/slog"
	"os"
	"strings"
)

type CtxKey string

const CtxName CtxKey = "logger"

type Oak struct {
	layers string
	logger *slog.Logger
}

func New(backend slog.Handler) *Oak {
	o := &Oak{
		logger: slog.New(backend),
	}
	return o
}

func NewDefault() *Oak {
	return New(TintHandler(os.Stdout, LevelTrace))
}

func (o *Oak) Trace(msg string, args ...any) {
	o.logger.Log(context.TODO(), LevelTrace, msg, o.attrs(args)...)
}

func (o *Oak) Debug(msg string, args ...any) {
	o.logger.Debug(msg, o.attrs(args)...)
}

func (o *Oak) Info(msg string, args ...any) {
	o.logger.Info(msg, o.attrs(args)...)
}

func (o *Oak) Notice(msg string, args ...any) {
	o.logger.Log(context.TODO(), LevelNotice, msg, o.attrs(args)...)
}

func (o *Oak) Warning(msg string, args ...any) {
	o.logger.Log(context.TODO(), LevelWarning, msg, o.attrs(args)...)
}

func (o *Oak) Error(msg string, args ...any) {
	o.logger.Error(msg, o.attrs(args)...)
}

func (o *Oak) Fatal(msg string, args ...any) {
	o.logger.Log(context.TODO(), LevelFatal, msg, o.attrs(args)...)
}

func (o *Oak) Layer(l string) *Oak {
	if l == "" {
		return o
	}
	if o.layers == "" {
		o.layers = l
	} else {
		ls := strings.Split(o.layers, ".")
		ls = append(ls, l)
		o.layers = strings.Join(ls, ".")
	}
	return o
}

func (o *Oak) PopLayer() *Oak {
	ls := strings.Split(o.layers, ".")
	ls = ls[:len(ls)-1]
	o.layers = strings.Join(ls, ".")
	return o
}

func (o *Oak) ClearLayers() *Oak {
	o.layers = ""
	return o
}

func (l *Oak) attrs(args []any) []any {
	if l.layers == "" {
		return args
	}
	return append(args, "layer", l.layers)
}

func (o *Oak) Handler() slog.Handler {
	return o.logger.Handler()
}
