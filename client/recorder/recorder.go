package recorder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"os"
	"time"
	"woole/connection/eventsource"
	"woole/console"
	"woole/util"

	"github.com/ecromaneli-golang/console/logger"
)

// StatusInternalProxyError is any unknown proxy error.
const StatusInternalProxyError = 999

var config = console.ReadConfig()
var log = logger.New("recorder")

var records = NewRecords(uint(config.MaxRecords))
var recorderHandler http.HandlerFunc
var proxyHandler http.HandlerFunc

func Start() {
	initializeTunnel()
}

func Retry(request *Request) {
	record := NewRecord(request)
	DoRequestAndStoreResponse(record)

	if log.IsDebugEnabled() {
		log.Debug(record.ToString())
	}
}

func GetRecords() *Records {
	return records
}

func initializeTunnel() {
	// Open connection with tunnel/request
	client, err := eventsource.NewRequest(config.TunnelURL() + "/request")
	if err != nil {
		log.Fatal("Failed to connect with tunnel on " + config.TunnelURL())
		os.Exit(1)
	}

	proxyHandler = createProxyHandler()

	// Receive events, parse data, do request, record them, and return response
	for event := range client.Stream {
		id := event.ID

		var req Request
		json.Unmarshal([]byte(event.Data.(string)), &req)

		go func() {
			record := NewRecordWithID(id, &req)
			DoRequestAndStoreResponse(record)
			sendResponseToServer(record)
		}()
	}
}

func DoRequestAndStoreResponse(record *Record) {
	res, elapsed := proxyRequest(record.Request)
	record.Response = res
	record.Elapsed = elapsed
	records.Add(record)
}

func sendResponseToServer(record *Record) {
	handleRedirections(record)

	resData, err := json.Marshal(*record.Response)
	panicIfNotNil(err)

	_, err = http.Post(config.TunnelURL()+"/response/"+record.ID, "application/json", bytes.NewBuffer(resData))
	panicIfNotNil(err)

	if log.IsDebugEnabled() {
		log.Debug(record.ToString())
	}
}

func handleRedirections(record *Record) {
	location := record.Response.Header.Get("location")
	if location != "" {
		record.Response.Header.Del("location")
		record.Response.Code = 200
		record.Response.Body = []byte("Trying to redirect to <a href='" + location + "'>" + location + "</a>...")
	}
}

func proxyRequest(req *Request) (*Response, time.Duration) {

	// Redirect and record the response
	recorder := httptest.NewRecorder()
	elapsed := util.Timer(func() {
		proxyHandler.ServeHTTP(recorder, req.ToHTTPRequest())
	})

	// Save req and res data
	return (&Response{}).FromResponseRecorder(recorder), elapsed

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
