package channel

import (
	"errors"
	"sync"

	"github.com/ecromaneli-golang/console/logger"
)

type Broker struct {
	mu        sync.RWMutex
	stopCh    chan any
	publishCh chan any
	subCh     chan chan any
	unsubCh   chan chan any
	running   bool
}

func NewBroker() *Broker {
	return &Broker{
		stopCh:    make(chan any),
		publishCh: make(chan any, 8),
		subCh:     make(chan chan any, 1),
		unsubCh:   make(chan chan any, 1),
	}
}

func (b *Broker) Start() {
	go b.start()
}

func (b *Broker) start() {
	b.running = true
	subs := map[chan any]any{}

	for {
		select {
		case <-b.stopCh:
			b.running = false
			for listener := range subs {
				close(listener)
			}
			return

		case listener := <-b.subCh:
			subs[listener] = nil

		case listener := <-b.unsubCh:
			delete(subs, listener)
			close(listener)

		case msg := <-b.publishCh:
			for sub := range subs {
				select {
				case sub <- msg:
				default:
					logger.GetInstance().Warn("The message was missed")
				}
			}
		}
	}
}

func (b *Broker) Subscribe() (chan any, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.running {
		listener := make(chan any, 8)
		b.subCh <- listener
		return listener, nil
	}

	return nil, errors.New("Trying to subscribe a stopped Broker")
}

func (b *Broker) Unsubscribe(listener chan any) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.running {
		b.unsubCh <- listener
	}
}

func (b *Broker) Publish(msg any) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if b.running {
		b.publishCh <- msg
		return nil
	}

	return errors.New("Trying to publish in a stopped Broker")
}

func (b *Broker) Stop() {
	b.mu.Lock()
	defer b.mu.Unlock()

	close(b.stopCh)
}

func (b *Broker) IsRunning() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return b.running
}
