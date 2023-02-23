package app

import (
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"

	"github.com/ecromaneli-golang/console/logger"
	"google.golang.org/grpc/credentials"
)

const (
	ClientToken = "{client}"
)

// Config has all the configuration parsed from the command line.
type Config struct {
	HostPattern     string
	HttpPort        string
	HttpsPort       string
	TunnelPort      string
	DashboardPort   string
	Timeout         int
	TlsCert         string
	TlsKey          string
	MaxRequestSize  int
	MaxResponseSize int
	mu              sync.Mutex
	isRead          bool
}

var config *Config = &Config{isRead: false}

func (cfg *Config) HasTlsFiles() bool {
	return cfg.TlsCert != "" && cfg.TlsKey != ""
}

// ReadConfig reads the arguments from the command line.
func ReadConfig() *Config {
	if !config.isRead {
		config.mu.Lock()
		defer config.mu.Unlock()
	}

	if config.isRead {
		return config
	}

	hostPattern := flag.String("pattern", "{client}", "Set the server host pattern. The '"+ClientToken+"' MUST be present to determine where to get client id")
	httpPort := flag.Int("http", 80, "HTTP Port")
	httpsPort := flag.Int("https", 443, "HTTPS Port")
	tunnelPort := flag.Int("tunnel", 8001, "Tunnel Port")
	dashboardPort := flag.Int("dashboard", 8000, "Dashboard Port")
	timeout := flag.Int("timeout", 10000, "Timeout for receive a response from Client")
	tlsCert := flag.String("tls-cert", "", "TLS cert/fullchain file path")
	tlsKey := flag.String("tls-key", "", "TLS key/privkey file path")
	logLevel := flag.String("log-level", "OFF", "Log Level")
	maxRequestSize := flag.Int("max-request-size", math.MaxInt32, "Maximum request size in bytes. 0 = max value")
	maxResponseSize := flag.Int("max-response-size", 4*1024*1024, "Maximum response size in bytes. 0 = max value")

	flag.Parse()

	logger.SetLogLevelStr(*logLevel)

	config = &Config{
		HostPattern:     *hostPattern,
		HttpPort:        strconv.Itoa(*httpPort),
		HttpsPort:       strconv.Itoa(*httpsPort),
		TunnelPort:      strconv.Itoa(*tunnelPort),
		DashboardPort:   strconv.Itoa(*dashboardPort),
		TlsCert:         *tlsCert,
		TlsKey:          *tlsKey,
		Timeout:         *timeout,
		MaxRequestSize:  *maxRequestSize,
		MaxResponseSize: *maxResponseSize,
		isRead:          true,
	}

	if !strings.Contains(config.HostPattern, ClientToken) {
		panic("Pattern MUST has " + ClientToken)
	}

	if config.MaxRequestSize == 0 {
		config.MaxRequestSize = math.MaxInt32
	}

	if config.MaxResponseSize == 0 {
		config.MaxResponseSize = math.MaxInt32
	}

	return config
}

func PrintConfig() {
	fmt.Println(ReadConfig())
}

func (cfg *Config) GetTransportCredentials() (credentials.TransportCredentials, error) {
	if !cfg.HasTlsFiles() {
		return nil, errors.New("TLS Files not provided")
	}

	// Load certificate and private key
	serverCert, err := tls.LoadX509KeyPair(cfg.TlsCert, cfg.TlsKey)
	if err != nil {
		return nil, err
	}

	// Create the credentials and return it
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.NoClientCert,
	}

	return credentials.NewTLS(tlsConfig), nil
}
