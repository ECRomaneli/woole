package app

import (
	"crypto/sha512"
	"crypto/tls"
	"flag"
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"
	"woole/internal/pkg/constants"
	"woole/internal/pkg/url"
	"woole/pkg/rand"

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
	serverKey              []byte
	available              bool
}

const (
	ClientToken = "{client}"
)

var (
	config   *Config = &Config{available: false}
	configMu sync.Mutex
)

func (cfg *Config) HasTlsFiles() bool {
	return cfg.TlsCert != "" && cfg.TlsKey != ""
}

// ReadConfig reads the arguments from the command line.
func ReadConfig() *Config {
	if !config.available {
		configMu.Lock()
		defer configMu.Unlock()
	}

	if config.available {
		return config
	}

	httpPort := flag.Int("http", url.GetDefaultPort("http"), "Port on which the server listens for HTTP requests")
	httpsPort := flag.Int("https", url.GetDefaultPort("https"), "Port on which the server listens for HTTPS requests")
	logLevel := flag.String("log-level", "OFF", "Level of detail for the logs to be displayed")
	hostnamePattern := flag.String("pattern", ClientToken, "Set the server hostname pattern. Example: "+ClientToken+".mysite.com to vary the subdomain")
	serverKey := flag.String("key", "", "Key used to hash the bearer")
	tlsCert := flag.String("tls-cert", "", "Path to the TLS certificate or fullchain file")
	tlsKey := flag.String("tls-key", "", "Path to the TLS private key file")
	tunnelPort := flag.Int("tunnel", constants.DefaultTunnelPort, "Port on which the gRPC tunnel listens")
	tunnelReconnectTimeout := flag.Int("tunnel-reconnect-timeout", 10000, "Timeout to reconnect the stream when lose connection")
	tunnelRequestSize := flag.Int("tunnel-request-size", math.MaxInt32, "Tunnel maximum request size in bytes")
	tunnelResponseSize := flag.Int("tunnel-response-size", math.MaxInt32, "Tunnel maximum response size in bytes")
	tunnelResponseTimeout := flag.Int("tunnel-response-timeout", 20000, "Timeout to receive a client response")

	flag.Parse()

	logger.SetLogLevelStr(*logLevel)

	config = &Config{
		HttpPort:               strconv.Itoa(*httpPort),
		HttpsPort:              strconv.Itoa(*httpsPort),
		HostnamePattern:        *hostnamePattern,
		serverKey:              []byte(*serverKey),
		TlsCert:                *tlsCert,
		TlsKey:                 *tlsKey,
		TunnelPort:             strconv.Itoa(*tunnelPort),
		TunnelReconnectTimeout: *tunnelReconnectTimeout,
		TunnelRequestSize:      *tunnelRequestSize,
		TunnelResponseSize:     *tunnelResponseSize,
		TunnelResponseTimeout:  *tunnelResponseTimeout,
		available:              true,
	}

	if !strings.Contains(config.HostnamePattern, ClientToken) {
		panic("Hostname pattern MUST has " + ClientToken)
	}

	if len(config.serverKey) == 0 {
		config.serverKey = rand.RandSha512("")
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

func GenerateBearer(clientKey []byte) []byte {
	hash := sha512.Sum512(append(config.serverKey, clientKey...))
	return hash[:]
}
