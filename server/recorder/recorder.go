package recorder

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"woole-server/console"
	"woole-server/util"

	"github.com/ecromaneli-golang/console/logger"
	"github.com/ecromaneli-golang/http/webserver"
)

// StatusInternalProxyError is any unknown proxy error.
const StatusInternalProxyError = 999

var config = console.ReadConfig()
var log = logger.New("recorder")

var records = NewRecords()

func ListenAndServe() error {
	server := webserver.NewServer()
	go serveTunnel()

	server.All("{client}.*.*/**", recorderHandler)
	return server.ListenAndServe(config.ServerPort)
}

func GetRecords() *Records {
	return records
}

func serveTunnel() {
	server := webserver.NewServer()

	server.WriteText("/", "<h1>Shh! We are listening here...</h1>")
	server.Get("{client}.*.*/request", requestSender)
	server.Post("{client}.*.*/response/{id}", responseReceiver)

	server.ListenAndServe(config.TunnelPort)
}

func recorderHandler(req *webserver.Request, res *webserver.Response) {
	client := validateClient(req.Param("client"), true)

	record := NewRecord((&Request{}).FromHTTPRequest(req.Raw))
	records.Add(client, record)

	record.Elapsed = util.Timer(func() {
		defer records.Remove(client, record.ID)

		select {
		case <-record.OnResponse.Receive():
		case <-time.After(time.Duration(config.Timeout) * time.Millisecond):
			webserver.NewHTTPError(http.StatusGatewayTimeout, "Record "+record.ID+" - Server timeout reached").Panic()
		case <-req.Raw.Context().Done():
			webserver.NewHTTPError(http.StatusGatewayTimeout, "Record "+record.ID+" - The request is no longer available").Panic()
		}
	})

	rec := record.Response

	// Write response
	res.Headers(rec.Header).Status(rec.Code).Write(rec.Body)

	if log.IsDebugEnabled() {
		log.Debug(client, "-", record.ToString())
	}
}

func requestSender(req *webserver.Request, res *webserver.Response) {
	client := validateClient(req.Param("client"), false)

	res.Headers(webserver.EventStreamHeader)

	log.Debug(client + " - Connection Established")
	defer log.Debug(client + " - Connection Finished")
	defer records.RemoveClient(client)

	records.clients[client] = NewRecordMap()

	tunnel := records.clients[client].Tunnel

	go func() {
		<-req.Raw.Context().Done()

		select {
		case tunnel <- nil:
		default:
		}
	}()

	for record := range tunnel {
		if req.IsDone() {
			return
		}

		res.FlushEvent(&webserver.Event{
			ID:   record.ID,
			Name: "request",
			Data: *record.Request,
		})
	}
}

func responseReceiver(req *webserver.Request, res *webserver.Response) {
	client := validateClient(req.Param("client"), true)

	record := records.FindByClientAndId(client, req.Param("id"))

	response := Response{}
	err := json.Unmarshal(req.Body(), &response)

	if err != nil {
		webserver.NewHTTPError(http.StatusBadRequest, err).Panic()
	}

	record.Response = &response
	record.OnResponse.SendLast()
}

func validateClient(client string, shouldExists bool) string {
	if len(client) == 0 {
		webserver.NewHTTPError(http.StatusForbidden, "The client provided no identification").Panic()
	}

	if records.ClientExists(client) != shouldExists {
		message := "The subdomain '" + client + "' is already in use"

		if shouldExists {
			message = "The subdomain '" + client + "' is not in use"
		}

		webserver.NewHTTPError(http.StatusServiceUnavailable, message).Panic()
	}

	return client
}

func splitHostPort(hostPort string) (host, port string) {
	host = hostPort

	colon := strings.LastIndexByte(host, ':')
	if colon == -1 {
		return host, ""
	}

	return hostPort[:colon], hostPort[colon+1:]
}
