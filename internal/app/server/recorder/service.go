package recorder

import (
	"net/http"
	"strings"
	"time"

	"woole/internal/app/server/app"
	"woole/internal/pkg/tunnel"
	"woole/pkg/timer"

	"woole/internal/app/server/recorder/adt"

	"github.com/ecromaneli-golang/http/webserver"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func getRecordWhenReady(client *adt.Client, req *webserver.Request) *adt.Record {
	record := adt.NewRecord((&tunnel.Request{}).FromHTTPRequest(req))
	record.Step = tunnel.Step_REQUEST
	client.AddRecord(record)

	var err error

	elapsed := timer.Exec(func() {
		defer client.RemoveRecord(record.Id)

		select {
		case <-record.OnResponse.Receive():
		case <-time.After(config.TunnelResponseTimeout):
			err = webserver.NewHTTPError(http.StatusGatewayTimeout, client.Id+" Record("+record.Id+") - Server timeout reached")
		case <-req.Raw.Context().Done():
			err = webserver.NewHTTPError(http.StatusGatewayTimeout, client.Id+" Record("+record.Id+") - The request is no longer available")
		}
	})

	if err != nil {
		record.Response = &tunnel.Response{Code: http.StatusGatewayTimeout, ServerElapsed: elapsed}
		logRecord(client.Id, record)
		panic(err)
	}

	record.Response.ServerElapsed = elapsed
	client.SendServerElapsed(record)

	return record
}

func sendServerMessage(stream tunnel.Tunnel_TunnelServer, client *adt.Client) {
	for record := range client.RecordChannel {
		err := stream.Send(&tunnel.ServerMessage{Record: record})

		if !handleGRPCErrors(err) {
			return
		}
	}
}

func receiveClientMessage(stream tunnel.Tunnel_TunnelServer, client *adt.Client) {
	for {
		tunnelRes, err := stream.Recv()

		if !handleGRPCErrors(err) {
			return
		}

		if tunnelRes.Record.Step != tunnel.Step_RESPONSE {
			log.Error("Wrong record step")
			return
		}

		if err == nil {
			client.SetRecordResponse(tunnelRes.Record.Id, tunnelRes.Record.Response)
		}

	}
}

func createSession(client *adt.Client, expireAt int64) *tunnel.Session {
	hostname := strings.Replace(config.HostnamePattern, app.ClientToken, client.Id, 1)

	auth := &tunnel.Session{
		ClientId:        client.Id,
		Hostname:        hostname,
		HttpPort:        config.HttpPort,
		ExpireAt:        expireAt,
		MaxRequestSize:  int32(config.TunnelRequestSize),
		MaxResponseSize: int32(config.TunnelResponseSize),
		ResponseTimeout: int64(config.TunnelResponseTimeout),
		Bearer:          client.Bearer,
	}

	if config.HasTlsFiles() {
		auth.HttpsPort = config.HttpsPort
	}

	return auth
}

func getClient(hs *tunnel.Handshake) (*adt.Client, error) {
	err := app.AuthClient(hs.PublicKey)
	if err != nil {
		log.Error(hs.ClientId, "-", err.Error())
		return nil, err
	}

	// Recover client session if exists
	client, err := clientManager.RecoverSession(hs.ClientId, hs.Bearer)

	if err != nil {
		log.Error(hs.ClientId, "-", err.Error())
		return nil, err
	}

	if client != nil {
		return client, nil
	}

	// Create session or try recover from other server with the same key
	client, err = clientManager.Register(hs.ClientId, hs.Bearer, app.GenerateBearer(hs.ClientKey))

	if err != nil {
		log.Error(hs.ClientId, "-", err.Error())
		return nil, err
	}

	log.Info(client.Id, "- Session Started")
	clientManager.DeregisterOnTimeout(client.Id, func() { log.Info(client.Id, "- Session Finished") })
	return client, nil
}

func logRecord(clientId string, record *adt.Record) {
	if log.IsInfoEnabled() {
		log.Info(clientId, "-", record.ToString(26))
	}
}

func panicIfNotNil(err any) {
	if err != nil {
		panic(err)
	}
}

// Handle gRPC errors and return if the error was or not handled
func handleGRPCErrors(err error) bool {
	if err == nil {
		return true
	}

	switch status.Code(err) {
	case codes.ResourceExhausted:
		log.Warn("Request discarded. Reason: Max size exceeded")
		return true
	case codes.Unavailable, codes.Canceled:
		return false
	default:
		log.Error(err)
		return false
	}
}
