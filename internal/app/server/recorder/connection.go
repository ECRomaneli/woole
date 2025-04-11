package recorder

import (
	"net"
	"net/http"
	"woole/internal/pkg/tunnel"

	web "woole/web/server"

	"github.com/ecromaneli-golang/http/webserver"
	"google.golang.org/grpc"
)

func serveWebServer() {
	server := webserver.NewServerWithFS(http.FS(web.EmbeddedFS))
	domain := config.GetDomain()

	if domain != "" {
		server.Get(domain+"/", recorderHandler)
	}

	server.All(config.HostnamePattern+"/**", recorderHandler)

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
	opts = append(opts, grpc.MaxRecvMsgSize(config.TunnelResponseSize))
	opts = append(opts, grpc.MaxSendMsgSize(config.TunnelRequestSize))

	if config.HasTlsFiles() {
		opts = append(opts, grpc.Creds(config.GetTransportCredentials()))
	}

	s := grpc.NewServer(opts...)

	tunnel.RegisterTunnelServer(s, &Tunnel{})

	go func() {
		if err := s.Serve(lis); err != nil {
			panic("Failed to serve Tunnel. Reason: " + err.Error())
		}
	}()
}
