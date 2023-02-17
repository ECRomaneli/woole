package dashboard

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"net/http"
	"strings"
	"woole/cmd/client/app"
	"woole/cmd/client/recorder"
	"woole/shared/payload"

	"github.com/ecromaneli-golang/http/webserver"
	"github.com/google/brotli/go/cbrotli"
)

//go:embed static
var embeddedFS embed.FS

var records = recorder.GetRecords()
var config = app.ReadConfig()

func ListenAndServe() error {
	return setupServer().ListenAndServe(":" + config.DashboardPort)
}

func setupServer() *webserver.Server {
	staticFolder, _ := fs.Sub(embeddedFS, "static")
	server := webserver.NewServerWithFS(http.FS(staticFolder))

	server.FileServer("/")
	server.Get("/record/stream", connHandler)
	server.Get("/record/{id}/response/body", responseBodyHandler)
	server.Get("/record/{id}/retry", retryHandler)
	server.Get("/record/clear", clearHandler)

	return server
}

func connHandler(req *webserver.Request, res *webserver.Response) {
	listener, err := records.Broker.Subscribe()
	panicIfNotNil(err)
	defer records.Broker.Unsubscribe(listener)

	res.Headers(webserver.EventStreamHeader)

	res.FlushEvent(&webserver.Event{
		Name: "info",
		Data: *(&Info{}).FromConfig(),
	})

	res.FlushEvent(&webserver.Event{
		Name: "records",
		Data: records.ThinClone(),
	})

	go func() {
		<-req.Raw.Context().Done()
		listener <- nil
	}()

	for msg := range listener {
		if msg == nil {
			break
		}

		res.FlushEvent(&webserver.Event{
			Name: "record",
			Data: msg.(*recorder.Record).ThinClone(),
		})
	}
}

func clearHandler(req *webserver.Request, res *webserver.Response) {
	records.RemoveAll()
	res.Status(http.StatusOK).NoBody()
}

func retryHandler(req *webserver.Request, res *webserver.Response) {
	record := records.FindById(req.Param("id"))
	recorder.Retry(record.Request)
}

func responseBodyHandler(req *webserver.Request, res *webserver.Response) {
	record := records.FindById(req.Param("id"))
	body := record.Response.Body
	res.WriteJSON(decompress(record.Response.GetHttpHeader().Get("Content-Encoding"), body))
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

	data, err = ioutil.ReadAll(reader)
	panicIfNotNil(err)

	return data
}

func dumpCurl(req *payload.Request) string {
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
