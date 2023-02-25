package recorder

import (
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"woole/cmd/server/app"
	"woole/shared/util"

	pb "woole/shared/payload"

	"github.com/ecromaneli-golang/console/logger"
	"github.com/ecromaneli-golang/http/webserver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	config        = app.ReadConfig()
	log           = logger.New("recorder")
	clientManager = NewClientManager()
)

type Tunnel struct{ pb.UnimplementedTunnelServer }

func Start() {
	serveTunnel()
	serveWebServer()
}

func serveWebServer() {
	server := webserver.NewServer()

	server.All(config.HostPattern+"/**", recorderHandler)

	if config.HasTlsFiles() {
		go func() {
			panic(server.ListenAndServeTLS(":"+config.HttpsPort, config.TlsCert, config.TlsKey))
		}()
	}

	panic(server.ListenAndServe(":" + config.HttpPort))
}

func serveTunnel() {
	lis, err := net.Listen("tcp", ":"+config.TunnelPort)
	panicIfNotNil(err)

	// Opts
	var opts []grpc.ServerOption
	opts = append(opts, grpc.MaxRecvMsgSize(config.MaxResponseSize))
	opts = append(opts, grpc.MaxSendMsgSize(config.MaxRequestSize))

	// TODO: Implement TLS CA on Client side
	// if config.HasTlsFiles() {
	// 	tlsCred, err := config.GetTransportCredentials()
	// 	if err == nil {
	// 		opts = append(opts, grpc.Creds(tlsCred))
	// 	} else {
	// 		log.Error("Failed to create Transport Credentials.", err)
	// 		panic(err)
	// 	}
	// }

	s := grpc.NewServer(opts...)

	pb.RegisterTunnelServer(s, &Tunnel{})

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatal("Failed to serve Tunnel. Reason: ", err)
			os.Exit(1)
		}
	}()
}

func (_t *Tunnel) Tunnel(stream pb.Tunnel_TunnelServer) error {
	// Receive the client ACK
	ack, err := stream.Recv()

	if !handleGRPCErrors(err) {
		return err
	}

	// Register the new client
	client := clientManager.Register(ack.GetClientId())

	// Schedule the client to deregister after the tunnel finishes
	log.Info(client.id + " - Connection Established")
	defer log.Info(client.id + " - Connection Finished")
	defer clientManager.Deregister(client.id)

	// Send the authentication data
	err = stream.Send(&pb.TunnelRequest{Auth: createAuth(client)})

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

func sendRequests(stream pb.Tunnel_TunnelServer, client *Client) {
	for record := range client.Tunnel {
		err := stream.Send(&pb.TunnelRequest{
			RecordId: record.Id,
			Request:  record.Request,
		})

		if !handleGRPCErrors(err) {
			return
		}
	}
}

func receiveResponses(stream pb.Tunnel_TunnelServer, client *Client) {
	for {
		tunnelRes, err := stream.Recv()

		if !handleGRPCErrors(err) {
			return
		}

		if err == nil {
			client.SetRecordResponse(tunnelRes.RecordId, tunnelRes.Response)
		}

	}
}

// handle gRPC errors and return if the error was or not handled
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

func createAuth(client *Client) *pb.Auth {
	url := strings.Replace(config.HostPattern, app.ClientToken, client.id, 1)

	auth := &pb.Auth{
		ClientId:        client.id,
		Url:             url,
		HttpPort:        config.HttpPort,
		TunnelPort:      config.TunnelPort,
		MaxRequestSize:  int32(config.MaxRequestSize),
		MaxResponseSize: int32(config.MaxResponseSize),
		Bearer:          string(client.bearer),
	}

	if config.HasTlsFiles() {
		auth.HttpsPort = config.HttpsPort
	}

	return auth
}

func recorderHandler(req *webserver.Request, res *webserver.Response) {
	clientId, err := hasClient(req.Param("client"))
	panicIfNotNil(err)

	client := clientManager.Get(clientId)

	record := NewRecord((&pb.Request{}).FromHTTPRequest(req))
	client.AddRecord(record)

	record.Elapsed = util.Timer(func() {
		defer client.RemoveRecord(record.Id)

		select {
		case <-record.OnResponse.Receive():
		case <-time.After(time.Duration(config.Timeout) * time.Millisecond):
			err = webserver.NewHTTPError(http.StatusGatewayTimeout, clientId+" Record("+record.Id+") - Server timeout reached")
		case <-req.Raw.Context().Done():
			err = webserver.NewHTTPError(http.StatusGatewayTimeout, clientId+" Record("+record.Id+") - The request is no longer available")
		}
	})

	if err != nil {
		record.Response = &pb.Response{Code: http.StatusGatewayTimeout}
		logRecord(clientId, record)
		panic(err)
	}

	// Write response
	recRes := record.Response
	res.Headers(recRes.GetHttpHeader()).Status(int(recRes.Code)).Write(recRes.Body)
	logRecord(clientId, record)
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

func logRecord(clientId string, record *Record) {
	if log.IsInfoEnabled() {
		log.Info(clientId, "-", record.ToString(26))
	}
}

func panicIfNotNil(err any) {
	if err != nil {
		panic(err)
	}
}
