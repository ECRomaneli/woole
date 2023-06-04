package adt

import (
	"fmt"
	"strconv"
	"sync"

	pb "woole/shared/payload"
	"woole/shared/util/channel"
	"woole/shared/util/sequence"
	"woole/shared/util/signal"
)

var seqId sequence.Seq

type Type string

const (
	DEFAULT  Type = "default"
	REPLAY   Type = "replay"
	REDIRECT Type = "redirect"
)

type Record struct {
	*pb.Record
	ClientId string `json:"clientId,omitempty"`
	Type     Type   `json:"type,omitempty"`
}

type Records struct {
	mu          sync.RWMutex
	records     map[string]*Record
	maxRecords  uint
	lastDeleted int
	signal      *signal.Signal
	Broker      *channel.Broker
}

func NewRecords(maxRecords uint) *Records {
	recs := &Records{
		records:     make(map[string]*Record),
		maxRecords:  maxRecords,
		signal:      signal.New(),
		Broker:      channel.NewBroker(),
		lastDeleted: 0,
	}
	recs.Broker.Start()
	return recs
}

func NewRecord(req *pb.Request, recType Type) *Record {
	id := seqId.NextString()
	return &Record{ClientId: id, Type: recType, Record: &pb.Record{Request: req}}
}

func EnhanceRecord(rec *pb.Record) *Record {
	return &Record{ClientId: seqId.NextString(), Record: rec, Type: DEFAULT}
}

func (recs *Records) AddRecordAndCallListeners(rec *Record) {
	recs.mu.Lock()
	defer recs.mu.Unlock()
	recs.records[rec.ClientId] = rec

	if len(recs.records) > int(recs.maxRecords) {
		recs.lastDeleted++
		delete(recs.records, strconv.Itoa(recs.lastDeleted))
	}

	recs.Broker.Publish(rec)
}

func (recs *Records) Get(id string) *Record {
	recs.mu.RLock()
	defer recs.mu.RUnlock()
	return recs.records[id]
}

func (recs *Records) GetByServerId(id string) *Record {
	recs.mu.RLock()
	defer recs.mu.RUnlock()

	for _, record := range recs.records {
		if record.Id == id {
			return record
		}
	}

	return nil
}

func (recs *Records) ResetServerIds() {
	recs.mu.Lock()
	defer recs.mu.Unlock()

	for _, record := range recs.records {
		record.Id = ""
	}
}

func (recs *Records) RemoveAll() {
	recs.mu.Lock()
	defer recs.mu.Unlock()

	recs.records = make(map[string]*Record)
	recs.signal.Send()
}

func (recs *Records) Each(iterator func(rec *Record)) {
	recs.mu.RLock()
	defer recs.mu.RUnlock()

	for _, rec := range recs.records {
		iterator(rec)
	}
}

func (recs *Records) ThinCloneWithoutResponseBody() *[]Record {
	slice := []Record{}

	recs.Each(func(r *Record) {
		slice = append(slice, *r.ThinCloneWithoutResponseBody())
	})

	return &slice
}

func (rec *Record) ThinCloneWithoutResponseBody() *Record {
	return &Record{
		ClientId: rec.ClientId,
		Type:     rec.Type,
		Record: &pb.Record{
			Id:      rec.Id,
			Request: rec.Request,
			Response: &pb.Response{
				Proto:   rec.Response.Proto,
				Status:  rec.Response.Status,
				Code:    rec.Response.Code,
				Header:  rec.Response.Header,
				Elapsed: rec.Response.Elapsed,
				/*Body: rec.Response.Body, Skipped */
			},
		},
	}
}

func (rec *Record) Clone() *Record {
	clone := &Record{
		ClientId: rec.ClientId,
		Type:     rec.Type,
		Record:   rec.Record.Clone(),
	}
	return clone
}

func (rec *Record) ToString(maxPathLength int) string {
	path := []byte(rec.Request.Path)

	if len(path) > maxPathLength {
		path = append([]byte("..."), path[len(path)-maxPathLength:]...)
	}

	method := "[" + rec.Request.Method + "]"

	strPathLength := strconv.Itoa(maxPathLength + 3)
	str := fmt.Sprintf("%8s %"+strPathLength+"s", method, string(path))

	if rec.Response == nil {
		return str
	}

	return str + fmt.Sprintf(" %d - %dms", rec.Response.Code, rec.Response.Elapsed)
}
