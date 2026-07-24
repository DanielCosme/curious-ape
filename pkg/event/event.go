package event

import "danicos.dev/daniel/curious-ape/pkg/core"

type Topic string

var (
	DayCreated  Topic = "day-created"
	HabitUpsert Topic = "habit-upsert"
)

type Event struct {
	Topic
	Date *core.Date

	// done is set by Publish so it can block until the subscriber finishes.
	// Subscribers must call Done() when handling is complete (prefer defer).
	done chan struct{}
}

// Done signals that this subscriber has finished handling the event.
// It is a no-op if the event was not delivered via a blocking Publish.
func (e Event) Done() {
	if e.done != nil {
		close(e.done)
	}
}

type EventChan chan Event
