package oak

import (
	"context"
	"io"
	"log/slog"
	"time"

	"github.com/lmittmann/tint"
)

const (
	LevelTrace   = slog.Level(-8) // Things like SQL queries, etc...
	LevelDebug   = slog.LevelDebug
	LevelInfo    = slog.LevelInfo
	LevelNotice  = slog.Level(2)
	LevelWarning = slog.LevelWarn
	LevelError   = slog.LevelError
	LevelFatal   = slog.Level(12) // Panic?
)

// NOTE: Slog handler guide
// https://github.com/golang/example/tree/master/slog-handler-guide
type QueuedHandler struct {
	queue        []slog.Record
	level        slog.Level
	traceStarted bool
	lastMsg      bool
}

func NewQueuedHandler(level slog.Level) *QueuedHandler {
	return &QueuedHandler{
		level: level,
	}
}

func (h *QueuedHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.level.Level()
}

// NOTE: The context is provided to support applications that provide logging information along the call chain.
// In a break with usual Go practice, the Handle method should not treat a canceled context as a signal to stop work.
func (h *QueuedHandler) Handle(_ context.Context, r slog.Record) error {
	if h.traceStarted {
		r.Message = "  " + r.Message
	} else if !h.lastMsg {
		h.traceStarted = true
	}
	h.queue = append(h.queue, r)
	return nil
}

func (h *QueuedHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// The logger used with this handler is not meant to be copied given that it queues it's records.
	panic("not implemented")
}

func (h *QueuedHandler) WithGroup(name string) slog.Handler {
	// The logger used with this handler is not meant to be copied given that it queues it's records.
	panic("not implemented")
}

func (h *QueuedHandler) EndTrace() {
	h.traceStarted = false
	h.lastMsg = true
}

func (h *QueuedHandler) Dequeue(backend slog.Handler) {
	for _, record := range h.queue {
		backend.Handle(context.Background(), record)
	}
	h.traceStarted = false
}

func TintHandler(out io.Writer, level slog.Leveler, withTimeKey bool) slog.Handler {
	return tint.NewHandler(out, &tint.Options{
		Level:      level,
		TimeFormat: time.StampMilli,
		NoColor:    false,
		AddSource:  false,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if withTimeKey == false && a.Key == slog.TimeKey && len(groups) == 0 {
				return slog.Attr{}
			}
			if a.Key == slog.LevelKey && len(groups) == 0 {
				level, ok := a.Value.Any().(slog.Level)
				if ok {
					switch level {
					case LevelTrace:
						return tint.Attr(7, slog.String(a.Key, "TRC")) // White
					case LevelDebug:
						return tint.Attr(5, slog.String(a.Key, "DBG")) // Magenta
					case LevelNotice:
						return tint.Attr(4, slog.String(a.Key, "NOT")) // Blue
					case LevelWarning:
						return tint.Attr(3, slog.String(a.Key, "WRN")) // Yellow
					case LevelError:
						return tint.Attr(1, slog.String(a.Key, "ERR")) // Red
					case LevelFatal:
						return tint.Attr(9, slog.String(a.Key, "FTL")) // Bright Red
					}
				}
			}
			return a
		},
	})
}
