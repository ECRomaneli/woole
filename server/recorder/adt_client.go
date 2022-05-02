package recorder

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"math/rand"
	"strconv"
	"sync"
	"time"
	"woole-server/util/sequence"
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
	}

	this.generateKey()
	return this
}

func (this *Client) NextId() string {
	return this.seq.NextString()
}

func (this *Client) Add(rec *Record) (id string) {
	this.mu.Lock()
	defer this.mu.Unlock()

	rec.Id = this.seq.NextString()

	this.data[rec.Id] = rec
	this.Tunnel <- rec

	return rec.Id
}

func (this *Client) Get(key string) *Record {
	this.mu.RLock()
	defer this.mu.RUnlock()

	return this.data[key]
}

func (this *Client) Remove(key string) *Record {
	this.mu.Lock()
	defer this.mu.Unlock()

	data := this.data[key]
	this.data[key] = nil

	return data
}

func (this *Client) Authorize(bearer string) bool {
	return bytes.Compare(this.bearer, []byte(bearer)[7:]) == 0
}

func (this *Client) generateKey() []byte {
	r1 := time.Now().UnixNano()
	r2 := rand.Int()
	r3 := rand.Int()

	h := sha1.New()
	h.Write([]byte(strconv.FormatInt(r1, 16) + this.name + strconv.Itoa(r2) + strconv.Itoa(r3)))
	this.bearer = []byte(hex.EncodeToString(h.Sum(nil)))

	return this.bearer
}
