package recorder

import (
	"context"
	"net/http"
	"time"
	"woole/internal/pkg/tunnel"

	"github.com/ecromaneli-golang/http/webserver"
)

// REST -> [ALL] /**
func recorderHandler(req *webserver.Request, res *webserver.Response) {
	clientId, err := hasClient(req.Param("client"))
	panicIfNotNil(err)

	client := clientManager.Get(clientId)
	if client.IsIdle {
		panic("Trying to use an idle client")
	}

	record := getRecordWhenReady(client, req)
	res.Headers(record.Response.GetHttpHeader()).Status(int(record.Response.Code)).Write(record.Response.Body)
	logRecord(clientId, record)
}

// RPC -> Tunnel(stream *TunnelServer)
func (_t *Tunnel) Tunnel(stream tunnel.Tunnel_TunnelServer) error {

	// Get the stream context
	ctx := stream.Context()

	// Receive the client handshake
	hs, err := stream.Recv()

	if !handleGRPCErrors(err) {
		return err
	}

	// Recover client session if exists
	client, err := getClient(hs.Handshake)
	if err != nil {
		return err
	}

	client.Connect()
	log.Info(client.Id, "- Tunnel Connected")

	var expireAt int64 = 0

	if config.TunnelConnectionTimeout != 0 {
		deadline := time.Now().Add(config.TunnelConnectionTimeout)
		expireAt = deadline.Unix()
		cancelableCtx, cancel := context.WithDeadline(stream.Context(), deadline)
		ctx = cancelableCtx
		defer cancel()
	}

	// Send session
	stream.Send(&tunnel.ServerMessage{Session: createSession(client, expireAt)})

	if !handleGRPCErrors(err) {
		return err
	}

	// Listen for HTTP responses from client
	go receiveClientMessage(stream, client)

	// Send new HTTP requests to client
	go sendServerMessage(stream, client)

	// Wait the end-of-stream
	<-ctx.Done()

	if ctx.Err() != context.DeadlineExceeded {
		log.Info(client.Id, "- Tunnel Disconnected")
		client.SetIdleTimeout(config.TunnelReconnectTimeout)
	} else {
		log.Info(client.Id, "- Tunnel Expired")
		client.SetIdleTimeout(0)
	}

	return ctx.Err()
}

// RPC -> TestConn()
func (_t *Tunnel) TestConn(_ context.Context, _ *tunnel.Empty) (*tunnel.Empty, error) {
	return new(tunnel.Empty), nil
}

func hasClient(clientId string) (string, error) {
	if len(clientId) == 0 {
		return clientId, webserver.NewHTTPError(http.StatusForbidden, "The client provided no identification")
	}

	if !clientManager.Exists(clientId) {
		message := "The client '" + clientId + "' is not in use"
		return clientId, webserver.NewHTTPError(http.StatusForbidden, message)
	}

	return clientId, nil
}
