package recorder

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"woole/internal/app/client/app"
	"woole/internal/app/client/recorder/adt"
	"woole/internal/pkg/tunnel"
	"woole/internal/pkg/url"
	"woole/pkg/timer"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Replay(request *tunnel.Request) {
	record := adt.NewRecord(request, adt.REPLAY)
	record.Response = proxyRequest(record.Request)
	records.AddRecordAndPublish(record)

	if log.IsInfoEnabled() {
		log.Info(record.ToString(26))
	}
}

func GetRecords() *adt.Records {
	return records
}

func onTunnelStart(client tunnel.TunnelClient, ctx context.Context, cancelCtx context.CancelFunc) error {
	defer cancelCtx()

	// Start the tunnel stream
	stream, err := client.Tunnel(ctx)
	if !handleGRPCErrors(err) {
		return err
	}

	// Send the handshake
	stream.Send(&tunnel.ClientMessage{Handshake: config.GetHandshake()})

	// Receive the session
	serverMsg, err := stream.Recv()
	if !handleGRPCErrors(err) {
		return err
	}

	app.SetSession(serverMsg.Session)

	// Reset old IDs
	records.ResetServerIds()

	// Listen for requests and send responses asynchronously
	for {
		serverMsg, err := stream.Recv()

		if err != nil {
			if !handleGRPCErrors(err) {
				return err
			}
			continue
		}

		go handleServerRecord(stream, serverMsg.Record)
	}
}

func handleServerRecord(stream tunnel.Tunnel_TunnelClient, serverRecord *tunnel.Record) {
	defer catchAllErrors(serverRecord)

	switch serverRecord.Step {
	case tunnel.Step_REQUEST:
		handleServerRequest(stream, serverRecord)
	case tunnel.Step_SERVER_ELAPSED:
		handleServerElapsed(stream, serverRecord)
	default:
		log.Error("Record Step Not Allowed")
	}
}

func handleServerRequest(stream tunnel.Tunnel_TunnelClient, serverRecord *tunnel.Record) {
	record := adt.EnhanceRecord(serverRecord)
	doRequest(record)

	err := stream.Send(&tunnel.ClientMessage{Record: record.ThinClone(tunnel.Step_RESPONSE)})
	if !handleGRPCErrors(err) {
		log.Error("Failed to send response for Record[", record.Id, "].", err)
	}

	records.AddRecordAndPublish(record)

	if log.IsInfoEnabled() {
		log.Info(record.ToString(26))
	}
}

func handleServerElapsed(stream tunnel.Tunnel_TunnelClient, serverRecord *tunnel.Record) {
	rec := records.GetByServerId(serverRecord.Id)

	if rec == nil {
		log.Warn("Record [", serverRecord.Id, "] is not available")
		return
	}

	rec.Step = tunnel.Step_SERVER_ELAPSED
	rec.Response.ServerElapsed = serverRecord.Response.ServerElapsed
	records.Publish(&adt.Record{ClientId: rec.ClientId, Record: serverRecord})
}

func doRequest(record *adt.Record) {
	record.Step = tunnel.Step_RESPONSE
	replaceUrlHeaderByCustomUrl(record.Request.Header, "Origin")
	replaceUrlHeaderByCustomUrl(record.Request.Header, "Referer")
	record.Response = proxyRequest(record.Request)
	handleRedirections(record)
}

func proxyRequest(req *tunnel.Request) *tunnel.Response {
	// Redirect and record the response
	recorder := httptest.NewRecorder()
	elapsed := timer.Exec(func() {
		proxyHandler.ServeHTTP(recorder, req.ToHTTPRequest())
	})

	// Save req and res data
	return (&tunnel.Response{}).FromResponseRecorder(recorder, elapsed)
}

func handleRedirections(record *adt.Record) {
	res := record.Response
	location := res.GetHeaderOrEmpty("Location")

	if location == "" {
		return
	}

	record.Type = adt.REDIRECT

	params := make(map[string]string)
	params["redirectUrl"] = location
	params["hostname"] = app.GetSessionWhenAvailable().Hostname

	newUrl, ok := url.ReplaceHostByUsingExampleStr(location, record.Request.Url)
	if !ok {
		panic("Error when trying to replace the host of [" + record.Request.Url + "]")
	}

	if newUrl.String() == record.Request.Url {
		params["enableCustomUrl"] = "false"
		params["customUrl"] = "#"
	} else {
		params["enableCustomUrl"] = "true"
		params["customUrl"] = newUrl.String()
	}

	record.Response.Body = []byte(app.RedirectTemplate.Apply(params))
	record.Response.Code = http.StatusOK

	res.SetHeader("Content-Type", "text/html")
	res.DelHeader("Location")
	res.DelHeader("Content-Encoding")
	res.SetHeader("Content-Length", strconv.Itoa(len(record.Response.Body)))
}

func replaceUrlHeaderByCustomUrl(header map[string]string, headerName string) {
	if header == nil {
		return
	}

	rawUrl := header[headerName]
	newUrl, ok := url.ReplaceHostByUsingExampleUrl(rawUrl, config.CustomUrl)

	if !ok {
		panic("Error when trying to replace the host of [" + rawUrl + "]")
	}

	header[headerName] = newUrl.String()
}

func handleGRPCErrors(err error) bool {
	if err == nil {
		return true
	}

	switch status.Code(err) {
	case codes.ResourceExhausted:
		log.Warn("Request discarded. Reason: Max size exceeded")
		return true
	default:
		return false
	}
}

func catchAllErrors(record *tunnel.Record) {
	err := recover()

	if err == nil {
		return
	}

	log.Error(err)
	// TODO: Improve error handling
}
