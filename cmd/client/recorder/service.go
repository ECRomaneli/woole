package recorder

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"woole/cmd/client/app"
	"woole/cmd/client/recorder/adt"
	pb "woole/shared/payload"
	"woole/shared/util"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Replay(request *pb.Request) {
	record := adt.NewRecord(request)
	DoRequest(record)

	if log.IsInfoEnabled() {
		log.Info(record.ToString(26))
	}
}

func GetRecords() *adt.Records {
	return records
}

func onTunnelStart(client pb.TunnelClient, ctx context.Context, cancelCtx context.CancelFunc) error {
	defer cancelCtx()

	// Start the tunnel stream
	stream, err := client.Tunnel(ctx)
	if !handleGRPCErrors(err) {
		return err
	}

	stream.Send(&pb.ClientMessage{Session: app.GetSession()})

	// Listen for requests and send responses asynchronously
	for {
		serverMsg, err := stream.Recv()

		if err != nil {
			if !handleGRPCErrors(err) {
				return err
			}
			continue
		}

		clientRecord := adt.EnhanceRecord(serverMsg.Record)

		if serverMsg.Record.Step == pb.Step_SEND_REQUEST {
			go handleServerRequest(stream, clientRecord)
		} else if serverMsg.Record.Step == pb.Step_SEND_SERVER_RESPONSE {
			go handleServerResponse(stream, clientRecord)
		} else {
			log.Error("Record Step Not Allowed")
		}
	}
}

func handleServerRequest(stream pb.Tunnel_TunnelClient, record *adt.Record) {
	record.Step = pb.Step_RECEIVE_RESPONSE
	replaceUrlHeaderByCustomUrl(record.Request.Header, "Origin")
	replaceUrlHeaderByCustomUrl(record.Request.Header, "Referer")

	DoRequest(record)
	handleRedirections(record)

	err := stream.Send(&pb.ClientMessage{
		Record: record.Record,
	})

	if log.IsInfoEnabled() {
		log.Info(record.ToString(26))
	}

	if !handleGRPCErrors(err) {
		log.Error("Failed to send response for Record[", record.Id, "].", err)
	}
}

func handleServerResponse(stream pb.Tunnel_TunnelClient, record *adt.Record) {
	log.Debug("Received a STEP.RECEIVE_SERVER_RESPONSE")
}

func DoRequest(record *adt.Record) {
	record.Response = proxyRequest(record.Request)
	records.Add(record)
}

func handleRedirections(record *adt.Record) {
	location := record.Response.GetHttpHeader().Get("location")
	if location != "" {
		httpHeader := record.Response.GetHttpHeader()
		httpHeader.Set("Content-Type", "text/html")
		httpHeader.Del("location")
		record.Response.Body = []byte("<!DOCTYPE html><html lang='en'><head><meta charset='utf-8'><title>Woole - Redirecting</title><meta name='viewport' content='width=device-width, initial-scale=1'></head><body><span>Trying to redirect to <a href='" + location + "'>" + location + "</a>...</span></body></html>")
		record.Response.Code = http.StatusOK
		httpHeader.Set("Content-Length", strconv.Itoa(len(record.Response.Body)))
		record.Response.SetHttpHeader(httpHeader)
	}
}

func proxyRequest(req *pb.Request) *pb.Response {
	// Redirect and record the response
	recorder := httptest.NewRecorder()
	elapsed := util.Timer(func() {
		proxyHandler.ServeHTTP(recorder, req.ToHTTPRequest())
	})

	// Save req and res data
	return (&pb.Response{}).FromResponseRecorder(recorder, elapsed)
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

func replaceUrlHeaderByCustomUrl(header map[string]*pb.StringList, headerName string) {
	if header == nil || header[headerName] == nil {
		return
	}

	referer := header[headerName].Val[0]

	refererUrl, err := url.Parse(referer)
	if err != nil {
		log.Error("Error when trying to parse [", referer, "] to URL. Reason: ", err.Error())
	}

	refererUrl.Scheme = config.CustomUrl.Scheme
	refererUrl.Host = config.CustomUrl.Host
	refererUrl.Opaque = config.CustomUrl.Opaque

	header[headerName] = &pb.StringList{Val: []string{refererUrl.String()}}
}
