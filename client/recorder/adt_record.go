package recorder

import (
	"fmt"
	"sync"
	"time"

	"woole/util/sequence"
	"woole/util/signal"
)

var seqId sequence.Seq

type Record struct {
	Id       string
	Request  *Request
	Response *Response
	Elapsed  time.Duration
}

type Records struct {
	sync.RWMutex
	records    []*Record
	maxRecords uint
	signal     *signal.Signal
}

func NewRecords(maxRecords uint) *Records {
	return &Records{maxRecords: maxRecords, signal: signal.New()}
}

func NewRecord(req *Request) *Record {
	return &Record{Id: seqId.NextString(), Request: req}
}

func NewRecordWithId(id string, req *Request) *Record {
	return &Record{Id: id, Request: req}
}

func (recs *Records) Add(rec *Record) {
	recs.Lock()
	defer recs.Unlock()

	recs.records = append(recs.records, rec)

	if len(recs.records) > int(recs.maxRecords) {
		recs.records = recs.records[1:]
	}

	recs.signal.Send()
}

func (recs *Records) FindById(id string) *Record {
	recs.RLock()
	defer recs.RUnlock()

	for _, record := range recs.records {
		if record.Id == id {
			return record
		}
	}

	return nil
}

func (recs *Records) GetLast() *Record {
	recs.RLock()
	defer recs.RUnlock()

	lastIndex := len(recs.records) - 1

	if lastIndex < 0 {
		return nil
	}

	return recs.records[lastIndex]
}

func (recs *Records) RemoveAll() {
	recs.Lock()
	defer recs.Unlock()

	recs.records = nil
	recs.signal.Send()
}

func (recs *Records) Each(iterator func(rec *Record)) {
	recs.RLock()
	defer recs.RUnlock()

	for _, rec := range recs.records {
		iterator(rec)
	}
}

func (recs *Records) Updated() bool {
	<-recs.signal.Receive()
	return true
}

func (recs *Records) OnUpdate(onUpdate func()) {
	recs.Updated()

	recs.RLock()
	defer recs.RUnlock()

	onUpdate()
}

func (this *Record) ToString() string {
	path := []byte(this.Request.Path)

	if len(path) > 25 {
		path = append([]byte("..."), path[len(path)-26:]...)
	}

	method := "[" + this.Request.Method + "]"

	str := fmt.Sprintf("%8s %30s", method, string(path))

	if this.Response == nil {
		return str
	}

	return str + fmt.Sprintf(" %d - %dms", this.Response.Code, this.Elapsed)
}
