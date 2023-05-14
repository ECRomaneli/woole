package recorder

import (
	"net/http"
	"strings"
	"time"

	"woole/cmd/server/app"

	"woole/cmd/server/recorder/adt"
	pb "woole/shared/payload"
	"woole/shared/util"

	"github.com/ecromaneli-golang/http/webserver"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func getRecordWhenReady(client *adt.Client, req *webserver.Request) *adt.Record {
	record := adt.NewRecord((&pb.Request{}).FromHTTPRequest(req))
	client.AddRecord(record)

	var err error

	elapsed := util.Timer(func() {
		defer client.RemoveRecord(record.Id)

		select {
		case <-record.OnResponse.Receive():
		case <-time.After(time.Duration(config.TunnelResponseTimeout) * time.Millisecond):
			err = webserver.NewHTTPError(http.StatusGatewayTimeout, client.Id+" Record("+record.Id+") - Server timeout reached")
		case <-req.Raw.Context().Done():
			err = webserver.NewHTTPError(http.StatusGatewayTimeout, client.Id+" Record("+record.Id+") - The request is no longer available")
		}
	})

	if err != nil {
		record.Response = &pb.Response{Code: http.StatusGatewayTimeout, Elapsed: elapsed}
		logRecord(client.Id, record)
		panic(err)
	}

	record.Response.Elapsed = elapsed
	return record
}

func sendRequests(stream pb.Tunnel_TunnelServer, client *adt.Client) {
	for record := range client.GetNewRecords() {
		err := stream.Send(&pb.ServerMessage{
			Record: &pb.Record{
				Id:      record.Id,
				Request: record.Request,
				Step:    pb.Step_SEND_REQUEST,
			},
		})

		if !handleGRPCErrors(err) {
			return
		}
	}
}

func receiveResponses(stream pb.Tunnel_TunnelServer, client *adt.Client) {
	for {
		tunnelRes, err := stream.Recv()

		if !handleGRPCErrors(err) {
			return
		}

		if tunnelRes.Record.Step != pb.Step_RECEIVE_RESPONSE {
			log.Error("Wrong record step")
			return
		}

		if err == nil {
			client.SetRecordResponse(tunnelRes.Record.Id, tunnelRes.Record.Response)
		}

	}
}

func createSession(client *adt.Client) *pb.Session {
	hostname := strings.Replace(config.HostnamePattern, app.ClientToken, client.Id, 1)

	auth := &pb.Session{
		ClientId:        client.Id,
		Hostname:        hostname,
		HttpPort:        config.HttpPort,
		MaxRequestSize:  int32(config.TunnelRequestSize),
		MaxResponseSize: int32(config.TunnelResponseSize),
		Bearer:          client.Bearer,
	}

	if config.HasTlsFiles() {
		auth.HttpsPort = config.HttpsPort
	}

	return auth
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
