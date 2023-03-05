package recorder

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"net/http/httputil"
	pb "woole/shared/payload"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func Start() {
	proxyHandler = createProxyHandler()
	startTunnelStream()
}

func createProxyHandler() http.HandlerFunc {
	proxy := httputil.NewSingleHostReverseProxy(config.CustomUrl)

	proxy.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, err error) {
		log.Error(err, ":", req.Method, req.URL)
		rw.WriteHeader(StatusInternalProxyError)
		fmt.Fprintf(rw, "%v", err)
	}

	return func(rw http.ResponseWriter, req *http.Request) {
		req.Host = config.CustomUrl.Host
		req.URL.Host = config.CustomUrl.Host
		req.URL.Scheme = config.CustomUrl.Scheme
		proxy.ServeHTTP(rw, req)
	}
}

func connectTunnel(enableTransportCredentials bool) (pb.Tunnel_TunnelClient, context.CancelFunc) {
	// Opts
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(math.MaxInt32)))
	opts = append(opts, grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(math.MaxInt32)))

	if enableTransportCredentials {
		opts = append(opts, grpc.WithTransportCredentials(config.GetTransportCredentials()))
	} else {
		log.Warn("Connecting with tunnel without TLS Credentials...")
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	// Dial tunnel
	conn, err := grpc.Dial(config.TunnelUrl.Host, opts...)
	exitIfNotNil("Failed to connect with tunnel on "+config.TunnelUrl.String(), err)

	// Create a cancelable context
	ctx, cancel := context.WithCancel(context.Background())

	// Start the tunnel stream
	client := pb.NewTunnelClient(conn)
	stream, err := client.Tunnel(ctx)

	// If unavailable, retry without credentials
	if status.Code(err) == codes.Unavailable && enableTransportCredentials {
		cancel()
		conn.Close()
		return connectTunnel(false)
	}

	exitIfNotNil("Failed to connect with tunnel on "+config.TunnelUrl.String(), err)

	return stream, func() {
		cancel()
		conn.Close()
	}
}
