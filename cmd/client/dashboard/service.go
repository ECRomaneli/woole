package dashboard

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"fmt"
	"io"
	"strings"
	"woole/cmd/client/app"
	"woole/cmd/client/recorder"
	pb "woole/shared/payload"

	"github.com/google/brotli/go/cbrotli"
)

var records = recorder.GetRecords()
var config = app.ReadConfig()

func ListenAndServe() error {
	return setupServer().ListenAndServe(":" + config.DashboardUrl.Port())
}

func decompress(contentEncoding string, data []byte) []byte {

	// "compress" content encoding is not supported yet
	if data == nil || contentEncoding == "" || contentEncoding == "identity" || contentEncoding == "compress" {
		return data
	}

	var reader io.ReadCloser
	var err error

	defer func() {
		if reader != nil {
			err = reader.Close()
			panicIfNotNil(err)
		}
	}()

	if contentEncoding == "gzip" {
		reader, err = gzip.NewReader(bytes.NewReader(data))
		panicIfNotNil(err)
	} else if contentEncoding == "br" {
		reader = cbrotli.NewReader(bytes.NewReader(data))
	} else if contentEncoding == "deflate" {
		reader = flate.NewReader(bytes.NewReader(data))
	}

	data, err = io.ReadAll(reader)
	panicIfNotNil(err)

	return data
}

func dumpCurl(req *pb.Request) string {
	var b strings.Builder
	// Build cmd.
	fmt.Fprintf(&b, "curl -X %s %s", req.Method, req.Url)
	// Build headers.
	for k, v := range req.GetHttpHeader() {
		fmt.Fprintf(&b, " \\\n  -H '%s: %s'", k, strings.Join(v, " "))
	}
	// Build body.
	if len(req.Body) > 0 {
		fmt.Fprintf(&b, " \\\n  -d '%s'", req.Body)
	}
	return b.String()
}

func panicIfNotNil(err any) {
	if err != nil {
		panic(err)
	}
}
