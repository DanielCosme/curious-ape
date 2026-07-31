package event

import (
	"danicos.dev/daniel/curious-ape/pkg/oak"
)

type Listener func(data any) error

type Bus interface {
	Publish(topic string, data any) error
	Subscribe(topic string, listener Listener)
}

type BusImpl struct {
	events map[string][]Listener
}

func NewBus() BusImpl {
	return BusImpl{
		events: map[string][]Listener{},
	}
}

func (b BusImpl) Publish(topic string, data any) error {
	if listeners, ok := b.events[topic]; ok {
		for _, run := range listeners {
			err := run(data)
			if err != nil {
				oak.Error("error executing event", "err", err.Error(), "topic", topic)
				return err
			}
		}
	} else {
		oak.Error("topic not found", "name", topic)
	}
	return nil
}

func (b BusImpl) Subscribe(topic string, listener Listener) {
	if listeners, found := b.events[topic]; found {
		listeners = append(listeners, listener)
	} else {
		b.events[topic] = []Listener{listener}
	}
}
