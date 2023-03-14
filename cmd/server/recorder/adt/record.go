package adt

import (
	"fmt"
	"strconv"
	"time"

	"woole/shared/payload"
	"woole/shared/util/signal"
)

type Record struct {
	Id         string
	Request    *payload.Request
	Response   *payload.Response
	Elapsed    time.Duration
	OnResponse *signal.Signal
}

func NewRecord(req *payload.Request) *Record {
	return &Record{
		Request:    req,
		OnResponse: signal.New(),
	}
}

func (recs *Record) ToString(maxPathLength int) string {
	path := []byte(recs.Request.Path)

	if len(path) > maxPathLength {
		path = append([]byte("..."), path[len(path)-maxPathLength:]...)
	}

	method := "[" + recs.Request.Method + "]"

	strPathLength := strconv.Itoa(maxPathLength + 3)
	str := fmt.Sprintf("%8s %"+strPathLength+"s", method, string(path))

	if recs.Response == nil {
		return str + fmt.Sprintf(" N/A - %dms", recs.Elapsed)
	}

	return str + fmt.Sprintf(" %3d - %dms", recs.Response.Code, recs.Elapsed)
}
