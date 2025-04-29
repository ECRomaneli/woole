package adt

import (
	"strconv"
	"strings"
	"sync"
	"time"

	"woole/internal/pkg/tunnel"
	"woole/pkg/channel"
	"woole/pkg/sequence"
	"woole/pkg/signal"
)

var seqId sequence.Seq

type Type string

const (
	DEFAULT  Type = "default"
	REPLAY   Type = "replay"
	REDIRECT Type = "redirect"
)

type Record struct {
	*tunnel.Record
	ClientId        string `json:"clientId,omitempty"`
	Type            Type   `json:"type,omitempty"`
	CreatedAtMillis int64  `json:"createdAtMillis,omitempty"`
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

func NewRecord(req *tunnel.Request, recType Type) *Record {
	id := seqId.NextString()
	createdAt := time.Now().UnixMilli()
	return &Record{ClientId: id, Type: recType, CreatedAtMillis: createdAt, Record: &tunnel.Record{Request: req}}
}

func EnhanceRecord(rec *tunnel.Record) *Record {
	createdAt := time.Now().UnixMilli()
	return &Record{ClientId: seqId.NextString(), CreatedAtMillis: createdAt, Record: rec, Type: DEFAULT}
}

func (recs *Records) AddRecord(rec *Record) {
	recs.mu.Lock()
	defer recs.mu.Unlock()
	recs.records[rec.ClientId] = rec

	if recs.maxRecords > 0 && len(recs.records) > int(recs.maxRecords) {
		recs.lastDeleted++
		delete(recs.records, strconv.Itoa(recs.lastDeleted))
	}
}

func (recs *Records) Publish(rec *Record) {
	recs.Broker.Publish(rec)
}

func (recs *Records) AddRecordAndPublish(rec *Record) {
	recs.AddRecord(rec)
	recs.Publish(rec)
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
		ClientId:        rec.ClientId,
		Type:            rec.Type,
		CreatedAtMillis: rec.CreatedAtMillis,
		Record: &tunnel.Record{
			Id:      rec.Id,
			Request: rec.Request,
			Response: &tunnel.Response{
				Proto:         rec.Response.Proto,
				Status:        rec.Response.Status,
				Code:          rec.Response.Code,
				Header:        rec.Response.Header,
				Elapsed:       rec.Response.Elapsed,
				ServerElapsed: rec.Response.ServerElapsed,
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

func (recs *Record) ToString(logRemoteAddr bool, maxPathLength int) string {
	var sb strings.Builder
	estSize := 9 + maxPathLength
	if logRemoteAddr {
		estSize += len(recs.Request.RemoteAddr) + 7
	}
	if recs.Response != nil {
		estSize += 18 // for elapsed time and status code
	}
	sb.Grow(estSize)

	// Write remote address if needed
	if logRemoteAddr && recs.Request.RemoteAddr != "" {
		sb.WriteString("From: ")
		sb.WriteString(recs.Request.RemoteAddr)
		sb.WriteByte(' ')
	}

	// Write method
	if len(recs.Request.Method) < 6 {
		sb.WriteString(strings.Repeat(" ", 6-len(recs.Request.Method)))
	}
	sb.WriteByte('[')
	sb.WriteString(recs.Request.Method)
	sb.WriteString("] ")

	// Write path
	path := recs.Request.Path
	if len(path) > maxPathLength {
		sb.WriteString("...")
		sb.WriteString(path[len(path)-maxPathLength:])
	} else {
		sb.WriteString(strings.Repeat(" ", maxPathLength+3-len(path)))
		sb.WriteString(path)
	}

	if recs.Response == nil {
		// Write N/A if response is nil
		sb.WriteString(" N/A")
		return sb.String()
	}

	// Write elapsed time
	sb.WriteString(" - ")
	sb.WriteString(strconv.FormatInt(int64(recs.Response.Elapsed), 10))
	sb.WriteString("ms")

	return sb.String()
}
