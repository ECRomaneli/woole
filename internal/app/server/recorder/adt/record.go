package adt

import (
	"fmt"
	"strconv"

	"woole/internal/pkg/tunnel"
	"woole/pkg/signal"
)

type Record struct {
	tunnel.Record
	OnResponse *signal.Signal
}

func NewRecord(req *tunnel.Request) *Record {
	return &Record{Record: tunnel.Record{Request: req}, OnResponse: signal.New()}
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
		return str + " N/A"
	}

	return str + fmt.Sprintf(" %3d - c: %dms / s: %dms", recs.Response.Code, recs.Response.Elapsed, recs.Response.ServerElapsed)
}
