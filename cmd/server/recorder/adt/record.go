package adt

import (
	"fmt"
	"strconv"

	pb "woole/shared/payload"
	"woole/shared/util/signal"
)

type Record struct {
	pb.Record
	OnResponse *signal.Signal
}

func NewRecord(req *pb.Request) *Record {
	return &Record{Record: pb.Record{Request: req}, OnResponse: signal.New()}
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
