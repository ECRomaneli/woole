package recorder

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"woole-server/app"
	"woole-server/util"

	"github.com/ecromaneli-golang/console/logger"
	"github.com/ecromaneli-golang/http/webserver"
)

var config = app.ReadConfig()
var log = logger.New("recorder")

var records = NewRecords()

func ListenAndServe() {
	server := webserver.NewServer()
	go serveTunnel()

	server.All(config.HostPattern+"/**", recorderHandler)

	if config.HasTlsFiles() {
		go func() {
			panic(server.ListenAndServeTLS(config.HttpsPort, config.TlsCert, config.TlsKey))
		}()
	}

	panic(server.ListenAndServe(config.HttpPort))
}

func GetRecords() *Records {
	return records
}

func serveTunnel() {
	server := webserver.NewServer()

	server.WriteText("/", "<h1>Shh! We are listening here...</h1>")
	server.Post("/register/{clientId}", registerClient)
	server.Get("/request/{clientId}", requestSender)
	server.Post("/response/{clientId}/{recordId}", responseReceiver)

	if !config.HasTlsFiles() {
		panic(server.ListenAndServe(config.TunnelPort))
	}

	panic(server.ListenAndServeTLS(config.TunnelPort, config.TlsCert, config.TlsKey))
}

func recorderHandler(req *webserver.Request, res *webserver.Response) {
	clientId := validateClient(req.Param("client"), true)
	panicIfClientLockNotMatch(clientId, true)

	record := NewRecord((&Request{}).FromHTTPRequest(req))
	records.Add(clientId, record)

	record.Elapsed = util.Timer(func() {
		defer records.Remove(clientId, record.Id)

		select {
		case <-record.OnResponse.Receive():
		case <-time.After(time.Duration(config.Timeout) * time.Millisecond):
			webserver.NewHTTPError(http.StatusGatewayTimeout, clientId+" ["+record.Id+"] - Server timeout reached").Panic()
		case <-req.Raw.Context().Done():
			webserver.NewHTTPError(http.StatusGatewayTimeout, clientId+" ["+record.Id+"] - The request is no longer available").Panic()
		}
	})

	rec := record.Response

	// Write response
	res.Headers(rec.Header).Status(rec.Code).Write(rec.Body)

	if log.IsDebugEnabled() {
		log.Debug(clientId, "-", record.ToString(26))
	}
}

func requestSender(req *webserver.Request, res *webserver.Response) {
	clientId := req.Param("clientId")

	client := validateAndAuthClient(clientId, req.Header("Authorization"))
	panicIfClientLockNotMatch(clientId, false)

	client.Lock()

	log.Trace(clientId + " - Connection Established")
	defer log.Trace(clientId + " - Connection Finished")
	defer records.RemoveClient(clientId)

	res.Headers(webserver.EventStreamHeader)

	go func() {
		<-req.Raw.Context().Done()

		select {
		case client.Tunnel <- nil:
		default:
		}
	}()

	for record := range client.Tunnel {
		if req.IsDone() {
			return
		}

		res.FlushEvent(&webserver.Event{
			ID:   record.Id,
			Name: "request",
			Data: *record.Request,
		})
	}
}

func responseReceiver(req *webserver.Request, res *webserver.Response) {
	client := validateAndAuthClient(req.Param("clientId"), req.Header("Authorization"))
	record := client.Get(req.Param("recordId"))

	if record == nil {
		return
	}

	response := Response{}
	err := json.Unmarshal(req.Body(), &response)

	if err != nil {
		webserver.NewHTTPError(http.StatusBadRequest, err).Panic()
	}

	record.Response = &response
	record.OnResponse.SendLast()
}

func registerClient(req *webserver.Request, res *webserver.Response) {
	clientId := req.Param("clientId")

	if records.ClientIsLocked(clientId) {
		count := 2
		for ; records.ClientIsLocked(clientId + strconv.Itoa(count)); count++ {
		}
		clientId = clientId + strconv.Itoa(count)
	}

	url := strings.Replace(config.HostPattern, app.ClientToken, clientId, 1)

	payload := app.AuthPayload{
		Name:   clientId,
		Http:   "http://" + url + config.HttpPort,
		Bearer: string(records.RegisterClient(clientId).bearer),
	}

	if config.HasTlsFiles() {
		payload.Https = "https://" + url + config.HttpsPort
	}

	res.WriteJSON(payload)
}

func validateAndAuthClient(clientId, bearer string) *Client {
	clientId = validateClient(clientId, true)

	client, err := records.Get(clientId, bearer)

	if err != nil {
		webserver.NewHTTPError(http.StatusUnauthorized, err).Panic()
	}

	return client
}

func panicIfClientLockNotMatch(clientId string, shouldBeLocked bool) {
	if len(clientId) == 0 {
		webserver.NewHTTPError(http.StatusForbidden, "The client provided no identification").Panic()
	}

	if records.ClientIsLocked(clientId) != shouldBeLocked {
		message := "The client '" + clientId + "' is already in use"

		if shouldBeLocked {
			message = "The client '" + clientId + "' is not in use"
		}

		webserver.NewHTTPError(http.StatusForbidden, message).Panic()
	}
}

func validateClient(clientId string, shouldExist bool) string {
	if len(clientId) == 0 {
		webserver.NewHTTPError(http.StatusForbidden, "The client provided no identification").Panic()
	}

	if records.ClientExists(clientId) != shouldExist {
		message := "The client '" + clientId + "' is already registered"

		if shouldExist {
			message = "The client '" + clientId + "' is not registered"
		}

		webserver.NewHTTPError(http.StatusForbidden, message).Panic()
	}

	return clientId
}

func splitHostPort(hostPort string) (host, port string) {
	host = hostPort

	colon := strings.LastIndexByte(host, ':')
	if colon == -1 {
		return host, ""
	}

	return hostPort[:colon], hostPort[colon+1:]
}
