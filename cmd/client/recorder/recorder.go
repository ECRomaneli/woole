package recorder

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"time"
	"woole/cmd/client/app"
	"woole/shared/payload"
	pb "woole/shared/payload"
	"woole/shared/util"

	"github.com/ecromaneli-golang/console/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const StatusInternalProxyError = 999

var config = app.ReadConfig()
var log = logger.New("recorder")

var records = NewRecords(uint(config.MaxRecords))
var proxyHandler http.HandlerFunc

func Replay(request *payload.Request) {
	record := NewRecord(request)
	DoRequest(record)

	if log.IsInfoEnabled() {
		log.Info(record.ToString(26))
	}
}

func GetRecords() *Records {
	return records
}

func startTunnelStream() {
	// Establish tunnel connection and retrieve request/response stream
	stream, cancelFunc := connectTunnel(true)
	defer cancelFunc()

	// Send ack with client id (if exists)
	stream.Send(&pb.TunnelResponse{ClientId: config.ClientId})

	// Retrieve the authentication data and store
	tunnelReq, err := stream.Recv()
	panicIfNotNil(err)

	if tunnelReq.Auth == nil {
		exitIfNotNil("Failed to authenticate with tunnel on "+config.TunnelUrl.String(), errors.New("authencation not sent"))
	}

	app.Authenticate(tunnelReq.Auth)

	// Listen for requests and send responses asynchronously
	for {
		tunnelReq, err := stream.Recv()

		if err != nil {
			if !handleGRPCErrors(err) {
				panic(err)
			}
			continue
		}

		go handleTunnelRequest(stream, tunnelReq)
	}
}

func handleTunnelRequest(stream pb.Tunnel_TunnelClient, tunnelReq *pb.TunnelRequest) {
	record := NewRecordWithId(tunnelReq.RecordId, tunnelReq.Request)
	replaceUrlHeaderByCustomUrl(record.Request.Header, "Origin")
	replaceUrlHeaderByCustomUrl(record.Request.Header, "Referer")
	DoRequest(record)
	handleRedirections(record)

	err := stream.Send(&pb.TunnelResponse{
		RecordId: record.Id,
		Response: record.Response,
	})

	if log.IsInfoEnabled() {
		log.Info(record.ToString(26))
	}

	if !handleGRPCErrors(err) {
		panic(err)
	}
}

func DoRequest(record *Record) {
	res, elapsed := proxyRequest(record.Request)
	record.Response = res
	record.Elapsed = elapsed
	records.Add(record)
}

func handleRedirections(record *Record) {
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

func proxyRequest(req *payload.Request) (*payload.Response, time.Duration) {
	// Redirect and record the response
	recorder := httptest.NewRecorder()
	elapsed := util.Timer(func() {
		proxyHandler.ServeHTTP(recorder, req.ToHTTPRequest())
	})

	// Save req and res data
	return (&payload.Response{}).FromResponseRecorder(recorder), elapsed
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

func panicIfNotNil(err error) {
	if err != nil {
		panic(err)
	}
}

func exitIfNotNil(msg string, err error) {
	if err != nil {
		fmt.Println(msg)
		log.Fatal(msg + ". Reason: " + err.Error())
		os.Exit(1)
	}
}
