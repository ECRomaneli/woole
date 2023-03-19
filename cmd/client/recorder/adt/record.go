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

type Record struct {
	pb.Record
	ClientId string
}

type Records struct {
	mu         sync.RWMutex
	records    []*Record
	maxRecords uint
	signal     *signal.Signal
	Broker     *channel.Broker
}

func NewRecords(maxRecords uint) *Records {
	records := &Records{maxRecords: maxRecords, signal: signal.New(), Broker: channel.NewBroker()}
	records.Broker.Start()
	return records
}

func NewRecord(req *pb.Request) *Record {
	id := seqId.NextString()
	return &Record{ClientId: id, Record: pb.Record{Id: id + "C", Request: req}}
}

func NewRecordWithId(id string, req *pb.Request) *Record {
	return &Record{ClientId: seqId.NextString(), Record: pb.Record{Id: id, Request: req}}
}

func (recs *Records) Add(rec *Record) {
	recs.mu.Lock()
	defer recs.mu.Unlock()

	recs.records = append(recs.records, rec)

	if len(recs.records) > int(recs.maxRecords) {
		recs.records = recs.records[1:]
	}

	recs.Broker.Publish(rec)
}

func (recs *Records) FindById(id string) *Record {
	recs.mu.RLock()
	defer recs.mu.RUnlock()

	for _, record := range recs.records {
		if record.Id == id {
			return record
		}
	}

	return nil
}

func (recs *Records) RemoveAll() {
	recs.mu.Lock()
	defer recs.mu.Unlock()

	recs.records = nil
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
		Record: pb.Record{
			Id: rec.Id,
			Request: &pb.Request{
				Proto:      rec.Request.Proto,
				Path:       rec.Request.Path,
				Method:     rec.Request.Method,
				RemoteAddr: rec.Request.RemoteAddr,
				Url:        rec.Request.Url,
				Header:     rec.Request.Header,
				Body:       rec.Request.Body,
			},
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
