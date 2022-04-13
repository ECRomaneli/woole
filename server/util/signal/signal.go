package signal

import "sync"

type Signal struct {
	rw         sync.RWMutex
	signalChan chan any
}

func New() *Signal {
	return &Signal{signalChan: make(chan any)}
}

func (this *Signal) Send() {
	this.rw.Lock()
	defer this.rw.Unlock()

	close(this.signalChan)
	this.signalChan = make(chan any)
}

func (this *Signal) SendLast() {
	this.rw.Lock()
	defer this.rw.Unlock()

	close(this.signalChan)
}

func (this *Signal) Receive() <-chan any {
	this.rw.RLock()
	defer this.rw.RUnlock()

	return this.signalChan
}
