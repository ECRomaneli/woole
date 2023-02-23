package recorder

import (
	"bytes"
	"sync"
	pb "woole/shared/payload"
	"woole/shared/util/hash"
	"woole/shared/util/sequence"
)

type Client struct {
	mu     sync.Mutex
	bearer []byte
	id     string
	seq    sequence.Seq
	data   map[string]*Record
	Tunnel chan *Record
}

func NewClient(id string) *Client {
	this := &Client{
		id:     id,
		Tunnel: make(chan *Record, 32),
		data:   make(map[string]*Record),
		bearer: hash.RandSha1(id),
	}
	return this
}

func (cl *Client) Authorize(bearer string) bool {
	return bytes.Equal(cl.bearer, []byte(bearer)[7:])
}

func (cl *Client) NextId() string {
	return cl.seq.NextString()
}

func (cl *Client) AddRecord(rec *Record) (id string) {
	rec.Id = cl.seq.NextString()
	cl.putRecord(rec.Id, rec)
	cl.Tunnel <- rec
	return rec.Id
}

func (cl *Client) RemoveRecord(recordId string) *Record {
	data := cl.data[recordId]
	cl.putRecord(recordId, nil)

	return data
}

func (cl *Client) SetRecordResponse(recordId string, response *pb.Response) {
	record := cl.data[recordId]

	if record == nil {
		return
	}

	record.Response = response
	record.OnResponse.SendLast()
}

func (cl *Client) putRecord(recordId string, record *Record) {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	cl.data[recordId] = record
}
