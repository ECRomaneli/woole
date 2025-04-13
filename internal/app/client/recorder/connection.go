package recorder

import (
	"context"
	"fmt"
	"math"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"time"
	"woole/internal/app/client/app"
	"woole/internal/pkg/tunnel"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func startConnectionWithServer(contextHandler func(tunnel.TunnelClient, context.Context, context.CancelFunc) (bool, error)) {
	firstConn := true
	var err error

	for attempt := 0; attempt <= config.MaxReconnectAttempts; attempt++ {
		if firstConn {
			firstConn = false
		} else {
			recoverOrExit(err)
		}

		// Establish tunnel connection and retrieve request/response stream
		client, ctx, cancelCtx, connErr := connectClient(config.EnableTLSTunnel)
		err = connErr

		if err != nil {
			continue
		}

		connEstablished, tunnelErr := contextHandler(client, ctx, cancelCtx)
		err = tunnelErr

		if connEstablished {
			attempt = 0
		}
	}
	recoverOrExit(status.Error(codes.Aborted, fmt.Sprintf("failed to establish connection after %d attempts", config.MaxReconnectAttempts)))
}

func CreateProxyHandler() http.HandlerFunc {
	proxy := httputil.NewSingleHostReverseProxy(config.CustomUrl)

	go setProxyTimeout()

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

func setProxyTimeout() {
	session := app.GetSessionWhenAvailable()

	// Customize the Transport to include a timeout
	http.DefaultTransport = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: time.Duration(session.ResponseTimeout),
		}).DialContext,
	}
}

func connectClient(enableTransportCredentials bool) (tunnel.TunnelClient, context.Context, context.CancelFunc, error) {
	// Opts
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(math.MaxInt32)))
	opts = append(opts, grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(math.MaxInt32)))

	if enableTransportCredentials {
		opts = append(opts, grpc.WithTransportCredentials(config.GetTransportCredentials()))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	// Dial server
	conn, err := grpc.NewClient(config.TunnelUrl.Host, opts...)
	if err != nil {
		conn.Close()
		return nil, nil, nil, err
	}

	// Create a cancelable context
	ctx, cancel := context.WithCancel(context.Background())
	cancelFn := func() { cancel(); conn.Close() }

	// Test connection and retry without credentials if needed
	client := tunnel.NewTunnelClient(conn)
	_, err = client.TestConn(ctx, new(tunnel.Empty))

	if err != nil {
		cancelFn()
		if status.Code(err) == codes.Unavailable && enableTransportCredentials {
			config.EnableTLSTunnel = false
			return connectClient(config.EnableTLSTunnel)
		}
		return nil, nil, nil, err
	}

	return client, ctx, cancelFn, nil
}

func recoverOrExit(err error) {
	errStatus, ok := status.FromError(err)

	if !ok || !isRecoverable(err) || !app.HasSession() {
		log.Fatal("[", config.TunnelUrl.String(), "]", errStatus.Code(), "-", errStatus.Message())
		log.Fatal("[", config.TunnelUrl.String(), "]", "Failed to connect with tunnel")
		os.Exit(1)
	}

	log.Error("[", config.TunnelUrl.String(), "]", errStatus.Code(), "-", errStatus.Message())
	if config.ReconnectInterval > 0 {
		log.Warn("[", config.TunnelUrl.String(), "]", "Trying to reconnect in", config.ReconnectIntervalStr, "...")
		<-time.After(config.ReconnectInterval)
	} else {
		log.Warn("[", config.TunnelUrl.String(), "]", "Trying to reconnect...")
	}
}

func isRecoverable(err error) bool {
	switch status.Code(err) {
	// e.g. server restart, load balance, etc
	case codes.Unavailable:
		return true
	case codes.Internal:
		return true
	default:
		return false
	}
}
