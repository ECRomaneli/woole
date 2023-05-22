package adt

import (
	"sync"
	"time"
	pb "woole/shared/payload"
	"woole/shared/util/sequence"
)

type Client struct {
	rw               sync.RWMutex
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

	if !client.Connect() {
		panic("Failed to connect client")
	}

	return client
}

func (cl *Client) AddRecord(rec *Record) (id string) {
	rec.Id = cl.seq.NextString()
	cl.putRecord(rec.Id, rec)
	cl.NewRecordChannel <- rec
	return rec.Id
}

func (cl *Client) RemoveRecord(recordId string) *Record {
	removedRecord := cl.getRecord(recordId)
	cl.putRecord(recordId, nil)
	return removedRecord
}

func (cl *Client) SetRecordResponse(recordId string, response *pb.Response) {
	record := cl.getRecord(recordId)

	if record == nil {
		return
	}

	record.Response = response
	record.OnResponse.SendLast()
}

func (cl *Client) DisconnectAfter(duration time.Duration) bool {
	return cl.IdleTimeout.Reset(duration)
}

func (cl *Client) Connect() bool {
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
