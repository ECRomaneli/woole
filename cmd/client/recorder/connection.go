package recorder

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"net/http/httputil"
	"os"
	"time"
	pb "woole/shared/payload"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func startConnectionWithServer() {
	for {
		// Establish tunnel connection and retrieve request/response stream
		client, ctx, cancelCtx, err := connectClient(config.EnableTLSTunnel)

		if err != nil {
			recoverOrExit(err)
			continue
		}

		err = onTunnelStart(client, ctx, cancelCtx)

		if err != nil {
			recoverOrExit(err)
		}
	}
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

func connectClient(enableTransportCredentials bool) (pb.TunnelClient, context.Context, context.CancelFunc, error) {
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
	conn, err := grpc.Dial(config.TunnelUrl.Host, opts...)
	if err != nil {
		conn.Close()
		return nil, nil, nil, err
	}

	// Create a cancelable context
	ctx, cancel := context.WithCancel(context.Background())
	cancelFn := func() { cancel(); conn.Close() }

	// Test connection and retry without credentials if needed
	client := pb.NewTunnelClient(conn)
	_, err = client.TestConn(ctx, new(pb.Empty))

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

	if ok && isRecoverable(err) {
		log.Error("[", config.TunnelUrl.String(), "]", errStatus.Code(), "-", errStatus.Message())
		log.Error("[", config.TunnelUrl.String(), "] Retrying in 5 seconds...")
		<-time.After(5 * time.Second)
	} else {
		log.Fatal("[", config.TunnelUrl.String(), "]", errStatus.Code(), "-", errStatus.Message())
		fmt.Println("Failed to connect with tunnel on " + config.TunnelUrl.String())
		os.Exit(1)
	}
}

func isRecoverable(err error) bool {
	switch status.Code(err) {
	// e.g. server restart, load balance, etc
	case codes.Unavailable:
		return true
	default:
		return false
	}
}
