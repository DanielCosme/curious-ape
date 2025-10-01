package oak

import (
	"log/slog"
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

func (o *Oak) Info(msg string, args ...any) {
	o.logger.Info(msg, o.attrs(args)...)
}

func (o *Oak) Error(msg string, args ...any) {
	o.logger.Error(msg, o.attrs(args)...)
}

func (o *Oak) Layer(l string) *Oak {
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
