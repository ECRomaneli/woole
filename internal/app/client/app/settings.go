package app

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"
	"sync"
	"time"
	"woole/internal/pkg/constants"
	"woole/internal/pkg/template"
	"woole/internal/pkg/tunnel"
	"woole/pkg/parser"
	"woole/pkg/rand"
	"woole/pkg/signal"
	web "woole/web/client"

	iurl "woole/internal/pkg/url"

	"github.com/ecromaneli-golang/console/logger"
	"google.golang.org/grpc/credentials"
)

// Config has all the configuration parsed from the command line.
type Config struct {
	ClientId             string
	ClientKey            []byte
	ProxyUrl             *url.URL
	HttpUrl              *url.URL
	TunnelUrl            *url.URL
	CustomUrl            *url.URL
	SnifferUrl           *url.URL
	MaxRecords           int
	ServerKey            string
	tlsSkipVerify        bool
	tlsCa                string
	EnableTLSTunnel      bool
	AllowReaders         bool
	IsStandalone         bool
	MaxReconnectAttempts int
	ReconnectIntervalStr string
	ReconnectInterval    time.Duration
	available            bool
}

const (
	defaultProxyPort         = "80"
	defaultSnifferPort       = "8000"
	defaultStandalonePort    = "8080"
	defaultStandaloneMessage = "[<hostname>]:<port>"
	defaultCustomUrlMessage  = "[<scheme>://]<hostname>[:<port>]"
)

var (
	RedirectTemplate                 = template.FromFile(web.EmbeddedFS, "redirect.html")
	config           *Config         = &Config{available: false}
	session          *tunnel.Session = &tunnel.Session{}
	sessionInitiated signal.Signal   = *signal.New()
	configMu         sync.Mutex
)

func HasSession() bool {
	return session.Bearer != nil
}

// If no session was provided yet, the routine will wait for a session
func GetSessionWhenAvailable() *tunnel.Session {
	<-sessionInitiated.Receive()
	return session
}

func SetSession(serverSession *tunnel.Session) {
	if !HasSession() {
		defer sessionInitiated.SendLast()
	}
	session = serverSession
}

func ReadConfig() *Config {
	if !config.available {
		configMu.Lock()
		defer configMu.Unlock()
	}

	if config.available {
		return config
	}

	emptyStr := ""

	clientId := flag.String("client", "", "Unique identifier of the client")
	httpUrl := flag.String("http", defaultStandaloneMessage, "Port to start the standalone server (disables tunnel)")
	proxyUrl := flag.String("proxy", defaultProxyPort, "URL of the target server to be proxied")
	tunnelUrl := flag.String("tunnel", constants.DefaultTunnelPortStr, "URL of the tunnel")
	customUrl := flag.String("custom-host", defaultCustomUrlMessage, "Custom host to be used when proxying")
	snifferPort := flag.String("sniffer", defaultSnifferPort, "Port on which the sniffing tool is available")
	maxRecords := flag.Int("records", 1000, "Max records to store. Use 0 to not define a limit.")
	logLevel := flag.String("log-level", "INFO", "Level of detail for the logs to be displayed")
	// allowReaders := flag.Bool("allow-readers", false, "Allow other connections to listen the requests")
	maxReconnectAttempts := flag.Int("reconnect-attempts", 5, "Maximum number of reconnection attempts. 0 for infinite")
	reconnectInterval := flag.String("reconnect-interval", "5s", "Time between reconnection attempts. Duration format")
	serverKey := flag.String("server-key", "", "Path to the ECC public key used to authenticate (only if configured by the server)")
	tlsSkipVerify := flag.Bool("tls-skip-verify", false, "Disables the validation of the integrity of the Server's certificate")
	tlsCa := flag.String("tls-ca", "", "Path to the TLS CA file. Only for self-signed certificates")

	flag.Usage = func() {
		flag.PrintDefaults()
		fmt.Printf("\n")
		fmt.Printf("Definitions:\n")
		fmt.Printf("  Duration Format\n")
		fmt.Printf("\tThe duration format allows you to specify time intervals using a combination of numeric values and time unit qualifiers up to \"(d)ay\":\n")
		fmt.Printf("\t- Example: \"1d 3h 50s\" for 1 day, 3 hours and 50 seconds.\n")
	}

	flag.Parse()

	logger.SetLogLevelStr(*logLevel)

	if *customUrl == defaultCustomUrlMessage {
		customUrl = proxyUrl
	}

	if *httpUrl == defaultStandaloneMessage {
		httpUrl = &emptyStr
	}

	config = &Config{
		ClientId:             *clientId,
		ClientKey:            rand.RandMD5(*clientId),
		HttpUrl:              iurl.RawUrlToUrl(*httpUrl, "http", defaultStandalonePort),
		ProxyUrl:             iurl.RawUrlToUrl(*proxyUrl, "http", ""),
		TunnelUrl:            iurl.RawUrlToUrl(*tunnelUrl, "grpc", constants.DefaultTunnelPortStr),
		CustomUrl:            iurl.RawUrlToUrl(*customUrl, "http", ""),
		SnifferUrl:           iurl.RawUrlToUrl(*snifferPort, "http", defaultSnifferPort),
		MaxRecords:           *maxRecords,
		ServerKey:            *serverKey,
		tlsSkipVerify:        *tlsSkipVerify,
		tlsCa:                *tlsCa,
		EnableTLSTunnel:      true,
		IsStandalone:         httpUrl != &emptyStr,
		MaxReconnectAttempts: *maxReconnectAttempts,
		ReconnectIntervalStr: *reconnectInterval,
		ReconnectInterval:    parseDurationOrPanic("reconnect-interval", *reconnectInterval),
		available:            true,
	}

	session.ClientId = *clientId
	return config
}

func (cfg *Config) GetTransportCredentials() credentials.TransportCredentials {
	tlsConfig := &tls.Config{InsecureSkipVerify: cfg.tlsSkipVerify}

	if tlsConfig.InsecureSkipVerify || cfg.tlsCa == "" {
		return credentials.NewTLS(tlsConfig)
	}

	caBytes, err := os.ReadFile(cfg.tlsCa)
	if err != nil {
		panic("Failed to load TLS CA File on: " + cfg.tlsCa + ". Reason: " + err.Error())
	}

	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM(caBytes) {
		panic("Failed to append certificate")
	}

	return credentials.NewTLS(tlsConfig)
}

func (cfg *Config) GetHandshake() *tunnel.Handshake {
	clientId := cfg.ClientId

	if HasSession() {
		clientId = session.ClientId
	}

	publicKey, err := loadPublicKeyECC(cfg.ServerKey)

	if err != nil {
		panic(err)
	}

	return &tunnel.Handshake{
		ClientId:  clientId,
		ClientKey: cfg.ClientKey,
		Bearer:    session.Bearer,
		PublicKey: publicKey,
	}
}

// loadPublicKeyECC loads the public key from the configured file
func loadPublicKeyECC(path string) ([]byte, error) {
	if path == "" {
		return nil, nil
	}

	keyData, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key file: %w", err)
	}

	var block *pem.Block
	for block == nil || block.Type != "PUBLIC KEY" {
		block, keyData = pem.Decode(keyData)
		if block == nil {
			return nil, errors.New("invalid public key format")
		}
	}
	return block.Bytes, nil
}

func ExpireAt() string {
	if session.ExpireAt == 0 {
		return "never"
	}

	t := time.Unix(session.ExpireAt, 0)
	return t.Format("2006-01-02 03:04:05 PM MST")
}

func parseDurationOrPanic(field string, duration string) time.Duration {
	dur, err := parser.ParseDuration(duration)
	if err != nil {
		panic(fmt.Sprintf("Invalid %s: %s", field, err))
	}
	return dur
}

func PrintConfig() {
	fmt.Println(ReadConfig())
}
