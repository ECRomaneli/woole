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
	"time"
	"woole/internal/pkg/constants"
	"woole/internal/pkg/url"
	"woole/pkg/parser"
	"woole/pkg/rand"

	"github.com/ecromaneli-golang/console/logger"
	"google.golang.org/grpc/credentials"
)

// Config has all the configuration parsed from the command line.
type Config struct {
	HostnamePattern         string
	HttpPort                string
	HttpsPort               string
	TlsCert                 string
	TlsKey                  string
	TunnelPort              string
	TunnelRequestSize       int
	TunnelResponseSize      int
	TunnelResponseTimeout   time.Duration
	TunnelReconnectTimeout  time.Duration
	TunnelConnectionTimeout time.Duration
	serverKey               []byte
	available               bool
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
	logLevel := flag.String("log-level", "INFO", "Level of detail for the logs to be displayed")
	hostnamePattern := flag.String("pattern", ClientToken, "Set the server hostname pattern. Example: "+ClientToken+".mysite.com to vary the subdomain")
	serverKey := flag.String("key", "", "Key used to hash the bearer")
	tlsCert := flag.String("tls-cert", "", "Path to the TLS certificate or fullchain file")
	tlsKey := flag.String("tls-key", "", "Path to the TLS private key file")
	tunnelPort := flag.Int("tunnel", constants.DefaultTunnelPort, "Port on which the gRPC tunnel listens")
	tunnelRequestSize := flag.String("tunnel-request-size", "", "Tunnel maximum request size. Size format (default \"2GB\", limited by gRPC)")
	tunnelResponseSize := flag.String("tunnel-response-size", "", "Tunnel maximum response size. Size format (default \"2GB\", limited by gRPC)")
	tunnelResponseTimeout := flag.String("tunnel-response-timeout", "10s", "Timeout to receive a client response. Duration format")
	tunnelReconnectTimeout := flag.String("tunnel-reconnect-timeout", "10s", "Timeout to reconnect the stream when lose connection. Duration format")
	tunnelConnectionTimeout := flag.String("tunnel-connection-timeout", "unset", "Timeout for client connections, Duration format")

	flag.Usage = func() {
		flag.PrintDefaults()
		fmt.Printf("\n")
		fmt.Printf("Definitions:\n")
		fmt.Printf("  Duration Format\n")
		fmt.Printf("\tThe duration format allows you to specify time intervals using a combination of numeric values and time unit qualifiers up to \"(d)ay\":\n")
		fmt.Printf("\t- Example: \"1d 3h 50s\" for 1 day, 3 hours and 50 seconds.\n")
		fmt.Printf("  Size Format\n")
		fmt.Printf("\tThe size format allows you to specify the number bytes using a combination of numeric values and unit qualifiers up to \"(T)era(B)ytes\":\n")
		fmt.Printf("\t- Example: \"30KB\" for 30 * 1024 Bytes.\n")
	}
	flag.Parse()

	logger.SetLogLevelStr(*logLevel)

	if *tunnelRequestSize == "" {
		tunnelRequestSize = strPointer("2gb")
	}
	if *tunnelResponseSize == "" {
		tunnelResponseSize = strPointer("2gb")
	}
	if *tunnelConnectionTimeout == "unset" {
		tunnelConnectionTimeout = strPointer("0")
	}

	config = &Config{
		HttpPort:                strconv.Itoa(*httpPort),
		HttpsPort:               strconv.Itoa(*httpsPort),
		HostnamePattern:         *hostnamePattern,
		serverKey:               []byte(*serverKey),
		TlsCert:                 *tlsCert,
		TlsKey:                  *tlsKey,
		TunnelPort:              strconv.Itoa(*tunnelPort),
		TunnelRequestSize:       parseTunnelMessageSizeOrPanic("tunnel-request-size", *tunnelRequestSize),
		TunnelResponseSize:      parseTunnelMessageSizeOrPanic("tunnel-response-size", *tunnelResponseSize),
		TunnelResponseTimeout:   parseDurationOrPanic("tunnel-response-timeout", *tunnelResponseTimeout),
		TunnelReconnectTimeout:  parseDurationOrPanic("tunnel-reconnect-timeout", *tunnelReconnectTimeout),
		TunnelConnectionTimeout: parseDurationOrPanic("tunnel-connection-timeout", *tunnelConnectionTimeout),
		available:               true,
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

func parseDurationOrPanic(field string, duration string) time.Duration {
	dur, err := parser.ParseDuration(duration)
	if err != nil {
		panic(fmt.Sprintf("Invalid %s: %s", field, err))
	}
	return dur
}

func parseTunnelMessageSizeOrPanic(field string, size string) int {
	sizeInt, err := parser.ParseBytes(size)
	if err != nil {
		panic(fmt.Sprintf("Invalid %s: %s", field, err))
	}

	if sizeInt == math.MaxInt32+1 {
		// "2GB - 1 byte" is the maximum size for gRPC
		sizeInt = math.MaxInt32
	} else if sizeInt > math.MaxInt32+1 {
		fmt.Printf("Warning: %s is greater than 2GB. Setting to 2GB", field)
	}

	return int(sizeInt)
}

func strPointer(str string) *string {
	return &str
}
