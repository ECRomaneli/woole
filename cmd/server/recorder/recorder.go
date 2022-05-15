package recorder

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"woole/cmd/server/app"
	"woole/shared/payload"
	"woole/shared/util"
	"woole/shared/util/hash"

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
			panic(server.ListenAndServeTLS(":"+config.HttpsPort, config.TlsCert, config.TlsKey))
		}()
	}

	panic(server.ListenAndServe(":" + config.HttpPort))
}

func GetRecords() *Records {
	return records
}

func serveTunnel() {
	server := webserver.NewServer()

	server.WriteText("/", "<h1>Shh! We are listening here...</h1>")
	server.Get("/request/{clientId?}", requestSender)
	server.Post("/response/{clientId}/{recordId}", responseReceiver)

	if !config.HasTlsFiles() {
		panic(server.ListenAndServe(":" + config.TunnelPort))
	}

	panic(server.ListenAndServeTLS(":"+config.TunnelPort, config.TlsCert, config.TlsKey))
}

func recorderHandler(req *webserver.Request, res *webserver.Response) {
	clientId, err := validateClient(req.Param("client"), true)
	panicIfNotNil(err)

	record := NewRecord((&payload.Request{}).FromHTTPRequest(req))
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

	if log.IsInfoEnabled() {
		log.Info(clientId, "-", record.ToString(26))
	}
}

func requestSender(req *webserver.Request, res *webserver.Response) {
	client, auth := registerClient(req.Param("clientId"))
	clientId := client.name

	log.Info(clientId + " - Connection Established")
	defer log.Info(clientId + " - Connection Finished")
	defer records.RemoveClient(clientId)

	res.Headers(webserver.EventStreamHeader)
	res.FlushEvent(&webserver.Event{Name: "auth", Data: auth})

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

	response := payload.Response{}
	err := json.Unmarshal(req.Body(), &response)

	if err != nil {
		webserver.NewHTTPError(http.StatusBadRequest, err).Panic()
	}

	record.Response = &response
	record.OnResponse.SendLast()
}

func registerClient(clientId string) (*Client, payload.Auth) {
	hasClientId := clientId != ""

	if !hasClientId {
		clientId = string(hash.RandSha1("")[:8])
	}

	clientId, err := validateClient(clientId, false)

	for err != nil {
		if hasClientId {
			clientId, err = validateClient(clientId+"-"+string(hash.RandSha1(clientId))[:5], false)
		} else {
			clientId, err = validateClient(string(hash.RandSha1(""))[:8], false)
		}
	}

	client := records.RegisterClient(clientId)
	url := strings.Replace(config.HostPattern, app.ClientToken, clientId, 1)

	auth := payload.Auth{
		ClientID:   clientId,
		URL:        url,
		HttpPort:   config.HttpPort,
		TunnelPort: config.TunnelPort,
		Bearer:     string(client.bearer),
	}

	if config.HasTlsFiles() {
		auth.HttpsPort = config.HttpsPort
	}

	return client, auth
}

func validateAndAuthClient(clientId, bearer string) *Client {
	clientId, err := validateClient(clientId, true)
	panicIfNotNil(err)

	client, err := records.Get(clientId, bearer)

	if err != nil {
		webserver.NewHTTPError(http.StatusUnauthorized, err).Panic()
	}

	return client
}

func validateClient(clientId string, shouldExist bool) (string, error) {
	if len(clientId) == 0 {
		return clientId, webserver.NewHTTPError(http.StatusForbidden, "The client provided no identification")
	}

	if records.ClientExists(clientId) != shouldExist {
		message := "The client '" + clientId + "' is already in use"

		if shouldExist {
			message = "The client '" + clientId + "' is not in use"
		}

		return clientId, webserver.NewHTTPError(http.StatusForbidden, message)
	}

	return clientId, nil
}

func splitHostPort(hostPort string) (host, port string) {
	host = hostPort

	colon := strings.LastIndexByte(host, ':')
	if colon == -1 {
		return host, ""
	}

	return hostPort[:colon], hostPort[colon+1:]
}

func panicIfNotNil(err any) {
	if err != nil {
		panic(err)
	}
}
