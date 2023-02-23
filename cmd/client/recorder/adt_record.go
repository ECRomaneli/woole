package recorder

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"woole/shared/payload"
	"woole/shared/util/channel"
	"woole/shared/util/sequence"
	"woole/shared/util/signal"
)

var seqId sequence.Seq

type Record struct {
	Id       string            `json:"id"`
	Request  *payload.Request  `json:"request"`
	Response *payload.Response `json:"response"`
	Elapsed  time.Duration     `json:"elapsed"`
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

func NewRecord(req *payload.Request) *Record {
	return &Record{Id: "R" + seqId.NextString(), Request: req}
}

func NewRecordWithId(id string, req *payload.Request) *Record {
	return &Record{Id: id, Request: req}
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

func (recs *Records) GetLast() *Record {
	recs.mu.RLock()
	defer recs.mu.RUnlock()

	lastIndex := len(recs.records) - 1

	if lastIndex < 0 {
		return nil
	}

	return recs.records[lastIndex]
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

func (recs *Records) Updated() bool {
	<-recs.signal.Receive()
	return true
}

func (recs *Records) ThinClone() *[]Record {
	slice := []Record{}

	recs.Each(func(r *Record) {
		slice = append(slice, *r.ThinClone())
	})

	return &slice
}

func (rec *Record) ThinClone() *Record {
	clone := &Record{Request: &payload.Request{}, Response: &payload.Response{}}

	clone.Id = rec.Id
	clone.Request.Url = rec.Request.Url
	clone.Request.Path = rec.Request.Path
	clone.Request.Method = rec.Request.Method
	clone.Request.Proto = rec.Request.Proto
	clone.Request.Header = rec.Request.Header
	clone.Request.Body = rec.Request.Body
	clone.Response.Code = rec.Response.Code
	clone.Response.Proto = rec.Response.Proto
	clone.Response.Header = rec.Response.Header
	clone.Elapsed = rec.Elapsed

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

	return str + fmt.Sprintf(" %d - %dms", rec.Response.Code, rec.Elapsed)
}
