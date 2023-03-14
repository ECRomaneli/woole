package adt

import (
	"sync"
	"time"
	pb "woole/shared/payload"
	"woole/shared/util/hash"
	"woole/shared/util/sequence"
)

type Client struct {
	mu            sync.Mutex
	Bearer        []byte
	Id            string
	seq           sequence.Seq
	records       map[string]*Record
	recordChannel chan *Record
	IdleTimeout   *time.Timer
}

func NewClient(id string) *Client {
	this := &Client{
		Id:            id,
		recordChannel: make(chan *Record, 32),
		records:       make(map[string]*Record),
		Bearer:        hash.RandSha512(id),
	}
	return this
}

func (cl *Client) GetNewRecords() chan *Record {
	return cl.recordChannel
}

func (cl *Client) AddRecord(rec *Record) (id string) {
	rec.Id = cl.seq.NextString()
	cl.putRecord(rec.Id, rec)
	cl.recordChannel <- rec
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

func (cl *Client) DisconnectAfter(duration time.Duration) {
	if cl.IdleTimeout == nil {
		cl.IdleTimeout = time.NewTimer(duration)
	} else {
		cl.IdleTimeout.Reset(duration)
	}
}

func (cl *Client) Connected() {
	if cl.IdleTimeout != nil {
		cl.IdleTimeout.Stop()
	}
}

func (cl *Client) putRecord(recordId string, record *Record) {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	cl.records[recordId] = record
}
