package recorder

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"woole-server/util/sequence"
	"woole-server/util/signal"
)

var seqID sequence.Seq

type RecordsByClient map[string]*RecordMap

type Record struct {
	ID         string
	Request    *Request
	Response   *Response
	Elapsed    time.Duration
	OnResponse *signal.Signal
}

type Records struct {
	sync.RWMutex
	clients RecordsByClient
}

func NewRecords() *Records {
	return &Records{clients: make(RecordsByClient)}
}

func NewRecord(req *Request) *Record {
	return &Record{
		ID:         seqID.NextString(),
		Request:    req,
		OnResponse: signal.New(),
	}
}

func (this *Records) ClientExists(client string) bool {
	return this.clients[client] != nil
}

func (this *Records) Add(client string, rec *Record) {
	this.Lock()
	defer this.Unlock()

	if this.clients[client] == nil {
		this.clients[client] = NewRecordMap()
	}

	recordMap := this.clients[client]

	recordMap.Put(rec.ID, rec)
}

func (this *Records) Remove(clientId, id string) *Record {
	this.Lock()
	defer this.Unlock()

	client := this.clients[clientId]

	if client == nil {
		return nil
	}

	return client.Remove(id)
}

func (this *Records) RemoveClient(client string) {
	this.Lock()
	defer this.Unlock()

	this.clients[client] = nil
}

func (this *Records) FindByClientAndId(client, id string) *Record {
	this.RLock()
	defer this.RUnlock()

	return this.clients[client].Get(id)
}

func (this *Record) ToString() string {
	path := []byte(this.Request.Path)

	if len(path) > 25 {
		path = append([]byte("..."), path[len(path)-26:]...)
	}

	method := "[" + this.Request.Method + "]"

	str := fmt.Sprintf("%8s %30s", method, string(path))

	if this.Response == nil {
		return str
	}

	return str + fmt.Sprintf(" %d - %dms", this.Response.Code, this.Elapsed)
}

// TODO Turn generic
type RecordMap struct {
	data   map[string]*Record
	keys   []string
	mu     *sync.RWMutex
	Tunnel chan *Record
	// onUpdate *signal.Signal
}

func NewRecordMap() *RecordMap {
	this := &RecordMap{mu: &sync.RWMutex{}, Tunnel: make(chan *Record, 16)}
	this.data = make(map[string]*Record)
	return this
}

func (this *RecordMap) Put(key string, value *Record) *RecordMap {
	this.mu.Lock()
	defer func() { this.Tunnel <- value }()
	defer this.mu.Unlock()

	this.data[key] = value

	return this
}

func (this *RecordMap) Get(key string) *Record {
	this.mu.RLock()
	defer this.mu.RUnlock()

	return this.data[key]
}

func (this *RecordMap) Remove(key string) *Record {
	this.mu.Lock()
	defer this.mu.Unlock()

	data := this.data[key]
	this.data[key] = nil

	return data
}

func (this *RecordMap) Last() *Record {
	this.mu.RLock()
	defer this.mu.RUnlock()

	length := this.Size()

	if length == 0 {
		return nil
	}

	return this.Get(this.keys[length-1])
}

func (this *RecordMap) Each(iterator func(*Record)) {
	this.mu.RLock()
	defer this.mu.RUnlock()

	for _, key := range this.keys {
		iterator(this.Get(key))
	}
}

func (this *RecordMap) Size() int {
	this.mu.RLock()
	defer this.mu.RUnlock()

	return len(this.keys)
}

// func (this *RecordMap) Updated() bool {
// 	<-this.onUpdate.Receive()
// 	return true
// }

// func (this *RecordMap) OnUpdate(onUpdate func()) {
// 	this.Updated()

// 	this.mu.RLock()
// 	defer this.mu.RUnlock()

// 	onUpdate()
// }

func (this *RecordMap) Shift() (*Record, error) {
	this.mu.Lock()
	defer this.mu.Unlock()

	if len(this.keys) == 0 {
		return nil, errors.New("Trying to shift an empty map")
	}

	i := this.data[this.keys[0]]
	delete(this.data, this.keys[0])
	this.keys = this.keys[1:]

	return i, nil
}
