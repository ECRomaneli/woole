package signal

import "sync"

type Signal struct {
	mu         sync.RWMutex
	signalChan chan any
}

func New() *Signal {
	return &Signal{signalChan: make(chan any)}
}

func (sig *Signal) Send() {
	sig.mu.Lock()
	defer sig.mu.Unlock()

	close(sig.signalChan)
	sig.signalChan = make(chan any)
}

func (sig *Signal) SendLast() {
	sig.mu.Lock()
	defer sig.mu.Unlock()

	close(sig.signalChan)
}

func (sig *Signal) Receive() <-chan any {
	sig.mu.RLock()
	defer sig.mu.RUnlock()

	return sig.signalChan
}
