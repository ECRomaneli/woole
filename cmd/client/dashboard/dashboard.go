package dashboard

import (
	"bytes"
	"compress/gzip"
	"embed"
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"woole/cmd/client/app"
	"woole/cmd/client/recorder"
	"woole/shared/payload"

	"github.com/ecromaneli-golang/http/webserver"
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
	server.Get("/conn", connHandler)
	server.Get("/info/{id}", infoHandler)
	server.Get("/clear", clearHandler)
	server.Get("/retry/{id}", retryHandler)

	return server
}

func connHandler(req *webserver.Request, res *webserver.Response) {
	res.Headers(webserver.EventStreamHeader)

	res.FlushEvent(&webserver.Event{
		Name: "config",
		Data: config,
	})

	res.FlushEvent(&webserver.Event{
		Name: "records",
		Data: (&Items{}).FromRecords(records),
	})

	var lastRecord *recorder.Record

	for !req.IsDone() {
		records.OnUpdate(func() { lastRecord = records.GetLast() })

		res.FlushEvent(&webserver.Event{
			Name: "record",
			Data: (&Item{}).FromRecord(lastRecord),
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

func infoHandler(req *webserver.Request, res *webserver.Response) {
	record := records.FindById(req.Param("id"))
	res.WriteJSON(dump(record))
}

func dump(rec *recorder.Record) map[string]string {
	req := rec.Request
	res := rec.Response
	return map[string]string{
		"request":  dumpContent(req.Header, req.Body, "%s %s %s\n\n", req.Method, req.Path, req.Proto),
		"response": dumpContent(res.Header, res.Body, "%s %s\n\n", res.Proto, res.Status),
		"curl":     dumpCurl(req),
	}
}

func dumpContent(header http.Header, body []byte, format string, args ...interface{}) string {
	b := strings.Builder{}
	fmt.Fprintf(&b, format, args...)
	dumpHeader(&b, header)
	b.WriteString("\n")
	dumpBody(&b, header, body)
	return b.String()
}

func dumpHeader(dst *strings.Builder, header http.Header) {
	var headers []string
	for k, v := range header {
		headers = append(headers, fmt.Sprintf("%s: %s\n", k, strings.Join(v, " ")))
	}
	sort.Strings(headers)
	for _, v := range headers {
		dst.WriteString(v)
	}
}

func dumpBody(dst *strings.Builder, header http.Header, body []byte) {
	reqBody := body
	if header.Get("Content-Encoding") == "gzip" {
		reader, _ := gzip.NewReader(bytes.NewReader(body))
		reqBody, _ = ioutil.ReadAll(reader)
	}
	dst.Write(reqBody)
}

func dumpCurl(req *payload.Request) string {
	var b strings.Builder
	// Build cmd.
	fmt.Fprintf(&b, "curl -X %s %s", req.Method, req.URL)
	// Build headers.
	for k, v := range req.Header {
		fmt.Fprintf(&b, " \\\n  -H '%s: %s'", k, strings.Join(v, " "))
	}
	// Build body.
	if len(req.Body) > 0 {
		fmt.Fprintf(&b, " \\\n  -d '%s'", req.Body)
	}
	return b.String()
}
