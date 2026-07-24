package event

import "sync"

type Broker interface {
	Publish(Event)
	Subscribe(Topic, EventChan)
}

var _ Broker = (*BrokerImpl)(nil)

type BrokerImpl struct {
	subs map[Topic][]EventChan
	mu   sync.Mutex
}

func NewBroker() *BrokerImpl {
	return &BrokerImpl{
		subs: map[Topic][]EventChan{},
		mu:   sync.Mutex{},
	}
}

// Publish delivers the event on each subscriber channel and blocks until that
// subscriber calls Event.Done(). Delivery is sequential per subscriber.
func (b *BrokerImpl) Publish(event Event) {
	b.mu.Lock()
	chans := append([]EventChan(nil), b.subs[event.Topic]...)
	b.mu.Unlock()

	for _, ch := range chans {
		ev := event
		ev.done = make(chan struct{})
		ch <- ev
		<-ev.done
	}
}

func (b *BrokerImpl) Subscribe(topic Topic, ch EventChan) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.subs[topic] = append(b.subs[topic], ch)
}
