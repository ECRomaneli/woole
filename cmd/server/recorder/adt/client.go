package adt

import (
	"sync"
	"time"
	pb "woole/shared/payload"
	"woole/shared/util/sequence"
)

type Client struct {
	mu               sync.Mutex
	Bearer           []byte
	Id               string
	seq              sequence.Seq
	records          map[string]*Record
	NewRecordChannel chan *Record
	IdleTimeout      *time.Timer
}

func NewClient(clientId string, bearer []byte) *Client {
	client := &Client{
		Id:               clientId,
		NewRecordChannel: make(chan *Record, 32),
		records:          make(map[string]*Record),
		Bearer:           bearer,
		IdleTimeout:      time.NewTimer(time.Minute),
	}
	client.Connected()
	return client
}

func (cl *Client) AddRecord(rec *Record) (id string) {
	rec.Id = cl.seq.NextString()
	cl.putRecord(rec.Id, rec)
	cl.NewRecordChannel <- rec
	return rec.Id
}

func (cl *Client) RemoveRecord(recordId string) *Record {
	data := cl.records[recordId]
	cl.putRecord(recordId, nil)

	return data
}

func (cl *Client) SetRecordResponse(recordId string, response *pb.Response) {
	record := cl.records[recordId]

	if record == nil {
		return
	}

	record.Response = response
	record.OnResponse.SendLast()
}

func (cl *Client) DisconnectAfter(duration time.Duration) bool {
	return cl.IdleTimeout.Reset(duration)
}

func (cl *Client) Connected() bool {
	return cl.IdleTimeout.Stop()
}

func (cl *Client) putRecord(recordId string, record *Record) {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	cl.records[recordId] = record
}
