package recorder

import (
	"context"
	"net/http"
	"time"
	pb "woole/shared/payload"

	"github.com/ecromaneli-golang/http/webserver"
)

// REST = [ALL] /**
func recorderHandler(req *webserver.Request, res *webserver.Response) {
	clientId, err := hasClient(req.Param("client"))
	panicIfNotNil(err)

	client := clientManager.Get(clientId)

	record := getRecordWhenReady(client, req)
	res.Headers(record.Response.GetHttpHeader()).Status(int(record.Response.Code)).Write(record.Response.Body)
	logRecord(clientId, record)
}

// gRPC = RequestSession(*Handshake) *Session
func (_t *Tunnel) RequestSession(_ctx context.Context, hs *pb.Handshake) (*pb.Session, error) {
	client := clientManager.Register(hs.ClientId)
	session := createSession(client)

	log.Info(client.Id, "- Session Started")

	go func() {
		<-client.IdleTimeout.C
		clientManager.Deregister(client.Id)
		log.Info(client.Id, "- Session Finished")
	}()

	client.DisconnectAfter(time.Duration(config.TunnelReconnectTimeout+1000) * time.Millisecond)
	return session, nil
}

// gRPC = Tunnel(stream *TunnelServer)
func (_t *Tunnel) Tunnel(stream pb.Tunnel_TunnelServer) error {
	// Receive the client ACK
	ack, err := stream.Recv()

	if !handleGRPCErrors(err) {
		return err
	}

	// Recover client session if exists
	client, err := clientManager.RecoverSession(ack.Session)
	if err != nil {
		log.Error(ack.Session.ClientId, "-", err.Error())
		return err
	}

	client.Connected()
	defer client.DisconnectAfter(time.Duration(config.TunnelReconnectTimeout) * time.Millisecond)

	// Schedule the client to deregister after the tunnel finishes
	log.Info(client.Id, "- Tunnel Connected")
	defer log.Info(client.Id, "- Tunnel Disconnected")

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
