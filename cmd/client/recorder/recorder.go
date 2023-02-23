package recorder

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"time"
	"woole/cmd/client/app"
	"woole/shared/payload"
	pb "woole/shared/payload"
	"woole/shared/util"

	"github.com/ecromaneli-golang/console/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

const StatusInternalProxyError = 999

var config = app.ReadConfig()
var log = logger.New("recorder")

var records = NewRecords(uint(config.MaxRecords))
var proxyHandler http.HandlerFunc

func Start() {
	proxyHandler = createProxyHandler()
	startTunnel()
}

func Retry(request *payload.Request) {
	record := NewRecord(request)
	DoRequest(record)

	if log.IsInfoEnabled() {
		log.Info(record.ToString(26))
	}
}

func GetRecords() *Records {
	return records
}

func startTunnel() {
	// Establish tunnel connection and retrieve request/response stream
	stream, cancelFunc := connectTunnel()
	defer cancelFunc()

	// Send ack with client id (if exists)
	stream.Send(&pb.TunnelResponse{ClientId: config.ClientId})

	// Retrieve the authentication data and store
	tunnelReq, err := stream.Recv()
	panicIfNotNil(err)

	if tunnelReq.Auth == nil {
		exitIfNotNil("Failed to authenticate with tunnel on "+config.TunnelHostPort(), errors.New("authencation not sent"))
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

func connectTunnel() (pb.Tunnel_TunnelClient, context.CancelFunc) {
	// Opts
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	opts = append(opts, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(math.MaxInt32)))
	opts = append(opts, grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(math.MaxInt32)))

	// Dial tunnel
	conn, err := grpc.Dial(config.TunnelHostPort(), opts...)
	exitIfNotNil("Failed to connect with tunnel on "+config.TunnelHostPort(), err)

	// Create a cancelable context
	ctx, cancel := context.WithCancel(context.Background())

	// Start the tunnel stream
	client := pb.NewTunnelClient(conn)
	stream, err := client.Tunnel(ctx)
	exitIfNotNil("Failed to connect with tunnel on "+config.TunnelHostPort(), err)

	return stream, func() {
		cancel()
		conn.Close()
	}
}

func handleTunnelRequest(stream pb.Tunnel_TunnelClient, tunnelReq *pb.TunnelRequest) {
	record := NewRecordWithId(tunnelReq.RecordId, tunnelReq.Request)
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
		record.Response.Body = []byte("<!doctype html><html><body>Trying to redirect to <a href='" + location + "'>" + location + "</a>...</body></html>")
		record.Response.Code = http.StatusOK
		httpHeader.Set("Content-Length", strconv.Itoa(len(record.Response.Body)))
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

func createProxyHandler() http.HandlerFunc {
	url, _ := url.Parse(config.CustomHost)
	proxy := httputil.NewSingleHostReverseProxy(url)

	proxy.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, err error) {
		log.Error(err, ":", req.Method, req.URL)
		rw.WriteHeader(StatusInternalProxyError)
		fmt.Fprintf(rw, "%v", err)
	}

	return proxy.ServeHTTP
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
