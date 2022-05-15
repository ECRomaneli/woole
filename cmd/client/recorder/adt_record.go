package recorder

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"woole/shared/payload"
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
	sync.RWMutex
	records    []*Record
	maxRecords uint
	signal     *signal.Signal
}

func NewRecords(maxRecords uint) *Records {
	return &Records{maxRecords: maxRecords, signal: signal.New()}
}

func NewRecord(req *payload.Request) *Record {
	return &Record{Id: seqId.NextString(), Request: req}
}

func NewRecordWithId(id string, req *payload.Request) *Record {
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

func (this *Record) ToString(maxPathLength int) string {
	path := []byte(this.Request.Path)

	if len(path) > maxPathLength {
		path = append([]byte("..."), path[len(path)-maxPathLength:]...)
	}

	method := "[" + this.Request.Method + "]"

	strPathLength := strconv.Itoa(maxPathLength + 3)
	str := fmt.Sprintf("%8s %"+strPathLength+"s", method, string(path))

	if this.Response == nil {
		return str
	}

	return str + fmt.Sprintf(" %d - %dms", this.Response.Code, this.Elapsed)
}

func (this *Record) ThinClone() *Record {
	clone := &Record{Request: &payload.Request{}, Response: &payload.Response{}}

	clone.Id = this.Id
	clone.Request.URL = this.Request.URL
	clone.Request.Path = this.Request.Path
	clone.Request.Method = this.Request.Method
	clone.Request.Proto = this.Request.Proto
	clone.Request.Header = this.Request.Header
	clone.Request.Body = this.Request.Body
	clone.Response.Code = this.Response.Code
	clone.Response.Proto = this.Response.Proto
	clone.Response.Header = this.Response.Header
	clone.Elapsed = this.Elapsed

	return clone
}

func (this *Records) ThinClone() *[]Record {
	slice := []Record{}

	this.Each(func(rec *Record) {
		slice = append(slice, *rec.ThinClone())
	})

	return &slice
}
