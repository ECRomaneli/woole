package recorder

import (
	"bytes"
	"sync"
	"woole/shared/util/hash"
	"woole/shared/util/sequence"
)

type Client struct {
	bearer []byte
	name   string
	seq    sequence.Seq
	data   map[string]*Record
	mu     *sync.RWMutex
	Tunnel chan *Record
}

func NewClient(name string) *Client {
	this := &Client{
		name:   name,
		mu:     &sync.RWMutex{},
		Tunnel: make(chan *Record, 32),
		data:   make(map[string]*Record),
		bearer: hash.RandSha1(name),
	}
	return this
}

func (cl *Client) NextId() string {
	return cl.seq.NextString()
}

func (cl *Client) Add(rec *Record) (id string) {
	cl.mu.Lock()
	defer cl.mu.Unlock()

	rec.Id = cl.seq.NextString()

	cl.data[rec.Id] = rec
	cl.Tunnel <- rec

	return rec.Id
}

func (cl *Client) Get(key string) *Record {
	cl.mu.RLock()
	defer cl.mu.RUnlock()

	return cl.data[key]
}

func (cl *Client) Remove(key string) *Record {
	cl.mu.Lock()
	defer cl.mu.Unlock()

	data := cl.data[key]
	cl.data[key] = nil

	return data
}

func (cl *Client) Authorize(bearer string) bool {
	return bytes.Equal(cl.bearer, []byte(bearer)[7:])
}
