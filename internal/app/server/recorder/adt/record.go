package adt

import (
	"strconv"
	"strings"

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
	sb.WriteString(" - c: ")
	sb.WriteString(strconv.FormatInt(int64(recs.Response.Elapsed), 10))
	sb.WriteString("ms / s: ")
	sb.WriteString(strconv.FormatInt(int64(recs.Response.ServerElapsed), 10))
	sb.WriteString("ms")

	return sb.String()
}
