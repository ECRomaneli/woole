package app

import (
	"crypto/tls"
	"flag"
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"
	"woole/shared/constants"
	"woole/shared/util"

	"github.com/ecromaneli-golang/console/logger"
	"google.golang.org/grpc/credentials"
)

// Config has all the configuration parsed from the command line.
type Config struct {
	HostnamePattern        string
	HttpPort               string
	HttpsPort              string
	TlsCert                string
	TlsKey                 string
	TunnelPort             string
	TunnelReconnectTimeout int
	TunnelRequestSize      int
	TunnelResponseSize     int
	TunnelResponseTimeout  int
	isRead                 bool
}

const (
	ClientToken = "{client}"
)

var (
	config        *Config = &Config{isRead: false}
	writingConfig sync.Mutex
)

func (cfg *Config) HasTlsFiles() bool {
	return cfg.TlsCert != "" && cfg.TlsKey != ""
}

// ReadConfig reads the arguments from the command line.
func ReadConfig() *Config {
	if !config.isRead {
		writingConfig.Lock()
		defer writingConfig.Unlock()
	}

	if config.isRead {
		return config
	}

	httpPort := flag.Int("http", util.GetDefaultPort("http"), "HTTP Port")
	httpsPort := flag.Int("https", util.GetDefaultPort("https"), "HTTPS Port")
	logLevel := flag.String("log-level", "OFF", "Log Level")
	hostnamePattern := flag.String("pattern", ClientToken, "Set the server hostname pattern. Example: Use "+ClientToken+".mysite.com to vary the subdomain as client ID")
	tlsCert := flag.String("tls-cert", "", "TLS cert/fullchain file path")
	tlsKey := flag.String("tls-key", "", "TLS key/privkey file path")
	tunnelPort := flag.Int("tunnel", constants.DefaultTunnelPort, "Tunnel Port")
	tunnelReconnectTimeout := flag.Int("tunnel-reconnect-timeout", 10000, "Timeout to reconnect the stream when lose connection")
	tunnelRequestSize := flag.Int("tunnel-request-size", math.MaxInt32, "Tunnel maximum request size in bytes. 0 = max value")
	tunnelResponseSize := flag.Int("tunnel-response-size", 4*1024*1024, "Tunnel maximum response size in bytes. 0 = max value")
	tunnelResponseTimeout := flag.Int("tunnel-response-timeout", 20000, "Timeout to receive a response from Client")

	flag.Parse()

	logger.SetLogLevelStr(*logLevel)

	config = &Config{
		HttpPort:               strconv.Itoa(*httpPort),
		HttpsPort:              strconv.Itoa(*httpsPort),
		HostnamePattern:        *hostnamePattern,
		TlsCert:                *tlsCert,
		TlsKey:                 *tlsKey,
		TunnelPort:             strconv.Itoa(*tunnelPort),
		TunnelReconnectTimeout: *tunnelReconnectTimeout,
		TunnelRequestSize:      *tunnelRequestSize,
		TunnelResponseSize:     *tunnelResponseSize,
		TunnelResponseTimeout:  *tunnelResponseTimeout,
		isRead:                 true,
	}

	if !strings.Contains(config.HostnamePattern, ClientToken) {
		panic("Hostname pattern MUST has " + ClientToken)
	}

	if config.TunnelRequestSize == 0 {
		config.TunnelRequestSize = math.MaxInt32
	}

	if config.TunnelResponseSize == 0 {
		config.TunnelResponseSize = math.MaxInt32
	}

	return config
}

func PrintConfig() {
	fmt.Println(ReadConfig())
}

func (cfg *Config) GetTransportCredentials() credentials.TransportCredentials {
	if !cfg.HasTlsFiles() {
		panic("TLS certificate and/or private key not provided")
	}

	// Load certificate and private key
	serverCert, err := tls.LoadX509KeyPair(cfg.TlsCert, cfg.TlsKey)
	if err != nil {
		panic("Failed to load TLS certificate and/or private key. Reason: " + err.Error())
	}

	// Create the credentials and return it
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.NoClientCert,
	}

	return credentials.NewTLS(tlsConfig)
}
