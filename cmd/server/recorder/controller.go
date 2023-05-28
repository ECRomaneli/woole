package recorder

import (
	"context"
	"net/http"
	"time"
	"woole/cmd/server/app"
	"woole/cmd/server/recorder/adt"
	pb "woole/shared/payload"

	"github.com/ecromaneli-golang/http/webserver"
)

// REST = [ALL] /**
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

// gRPC = Tunnel(stream *TunnelServer)
func (_t *Tunnel) Tunnel(stream pb.Tunnel_TunnelServer) error {
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

	defer client.DisconnectAfter(time.Duration(config.TunnelReconnectTimeout) * time.Millisecond)
	defer log.Info(client.Id, "- Tunnel Disconnected")

	// Send session
	stream.Send(&pb.ServerMessage{Session: createSession(client)})

	if !handleGRPCErrors(err) {
		return err
	}

	// Listen for HTTP responses from client
	go receiveResponses(stream, client)

	// Send new HTTP requests to client
	go sendRequests(stream, client)

	// Wait the end-of-stream
	<-stream.Context().Done()
	return nil
}

// gRPC = TestConn()
func (_t *Tunnel) TestConn(_ context.Context, _ *pb.Empty) (*pb.Empty, error) {
	return new(pb.Empty), nil
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

func getClient(hs *pb.Handshake) (*adt.Client, error) {
	// Recover client session if exists
	client, err := clientManager.RecoverSession(hs.ClientId, hs.Bearer)

	if err != nil {
		log.Error(hs.ClientId, "-", err.Error())
		return nil, err
	}

	if client != nil {
		return client, nil
	}

	// Create session
	client = clientManager.Register(hs.ClientId, app.GenerateBearer(hs.ClientKey))

	if len(hs.Bearer) != 0 {
		// Verify if old session is equal to the new one
		client, err = clientManager.RecoverSession(hs.ClientId, hs.Bearer)

		if err != nil {
			clientManager.Deregister(hs.ClientId)
			log.Error(hs.ClientId, "-", err.Error())
			return nil, err
		}
	}

	log.Info(client.Id, "- Session Started")
	clientManager.DeregisterIfIdle(client.Id, func() { log.Info(client.Id, "- Session Finished") })
	return client, nil
}
