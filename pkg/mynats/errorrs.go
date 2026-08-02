package mynats

import (
	"context"
	"errors"
	"log/slog"

	"github.com/nats-io/nats.go"
)

func ErrNatsIsNoResponders(err error) bool {
	if err != nil {
		switch {
		case errors.Is(err, nats.ErrTimeout):
			// Request timed out – no reply received in time
			slog.Info("request timed out")

		case errors.Is(err, nats.ErrNoResponders):
			// No service is currently listening on that subject
			slog.Info("IGNORING: no responders available error")
			return true
		case errors.Is(err, nats.ErrConnectionClosed):
			slog.Info("connection is closed")

		case errors.Is(err, nats.ErrConnectionDraining):
			slog.Info("connection is draining")

		case errors.Is(err, context.DeadlineExceeded):
			// Only relevant when using RequestWithContext
			slog.Info("context deadline exceeded")

		case errors.Is(err, context.Canceled):
			slog.Info("context canceled")

		default:
			// Any other error
			slog.Error("request failed: " + err.Error())
		}
		return false
	}
	return false
}
