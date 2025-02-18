package pubsub

import (
	"sync"

	"github.com/pkg/errors"
)

var ErrTopicNotFound = errors.New("topic not found")

type Publisher struct {
	subs sync.Map
}

func NewPublisher() *Publisher {
	return &Publisher{}
}

func (ps *Publisher) Subscribe(topic string) chan any {
	ch := make(chan any, 1)
	subscribers, _ := ps.subs.LoadOrStore(topic, make([]chan any, 0))
	ps.subs.Store(topic, append(subscribers.([]chan any), ch))
	return ch
}

func (ps *Publisher) Publish(topic string, msg any) error {
	value, ok := ps.subs.Load(topic)
	if ok {
		subscribers := value.([]chan any)
		for _, ch := range subscribers {
			go func(ch chan any) {
				ch <- msg
			}(ch)
		}
	} else {
		return errors.Wrap(ErrTopicNotFound, "Publish")
	}

	return nil
}

func (ps *Publisher) Unsubscribe(topic string, ch chan any) {
	value, ok := ps.subs.Load(topic)
	if ok {
		subscribers := value.([]chan any)
		for i, subscriber := range subscribers {
			if subscriber == ch {
				subscribers = append(subscribers[:i], subscribers[i+1:]...)
				ps.subs.Store(topic, subscribers)
				close(ch)
				return
			}
		}
	}
}
