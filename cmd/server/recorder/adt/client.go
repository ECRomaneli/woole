package adt

import (
	"sync"
	"time"
	pb "woole/internal/pkg/payload"
	"woole/pkg/sequence"
)

type Client struct {
	rw            sync.RWMutex
	Bearer        []byte
	Id            string
	seq           sequence.Seq
	records       map[string]*Record
	RecordChannel chan *pb.Record
	IdleTimeout   *time.Timer
	IsIdle        bool
}

func NewClient(clientId string, bearer []byte) *Client {
	client := &Client{
		Id:            clientId,
		RecordChannel: make(chan *pb.Record, 32),
		records:       make(map[string]*Record),
		Bearer:        bearer,
		IdleTimeout:   time.NewTimer(time.Minute),
	}

	if !client.Connect() {
		panic("Failed to connect client")
	}

	return client
}

func (cl *Client) AddRecord(rec *Record) (id string) {
	rec.Id = cl.seq.NextString()
	cl.putRecord(rec.Id, rec)
	cl.RecordChannel <- &rec.Record
	return rec.Id
}

func (cl *Client) RemoveRecord(recordId string) *Record {
	removedRecord := cl.getRecord(recordId)
	cl.putRecord(recordId, nil)
	return removedRecord
}

func (cl *Client) SendServerElapsed(rec *Record) {
	cl.RecordChannel <- rec.ThinClone(pb.Step_SERVER_ELAPSED)
}

func (cl *Client) SetRecordResponse(recordId string, response *pb.Response) {
	record := cl.getRecord(recordId)

	if record == nil {
		return
	}
	record.Step = pb.Step_RESPONSE
	record.Response = response
	record.OnResponse.SendLast()
}

func (cl *Client) DisconnectAfter(duration time.Duration) bool {
	cl.IsIdle = true
	return cl.IdleTimeout.Reset(duration)
}

func (cl *Client) Connect() bool {
	cl.IsIdle = false
	return cl.IdleTimeout.Stop()
}

func (cl *Client) getRecord(recordId string) *Record {
	cl.rw.RLock()
	defer cl.rw.RUnlock()
	return cl.records[recordId]
}

func (cl *Client) putRecord(recordId string, record *Record) {
	cl.rw.Lock()
	defer cl.rw.Unlock()
	cl.records[recordId] = record
}
