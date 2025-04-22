package adt

import (
	"testing"

	"woole/internal/pkg/tunnel"
)

func makeRecord() *Record {
	req := &tunnel.Request{
		Path:       "/some/very/long/path/to/resource/with/query?foo=bar&baz=qux",
		Method:     "GET",
		RemoteAddr: "192.0.2.1:12345",
	}
	rec := NewRecord(req, DEFAULT)
	rec.Response = &tunnel.Response{
		Code:          200,
		Elapsed:       123,
		ServerElapsed: 45,
	}
	return rec
}

func BenchmarkToString(b *testing.B) {
	rec := makeRecord()
	for i := 0; i < b.N; i++ {
		_ = rec.ToString(true, 32)
	}
}
