package app

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"net/url"
	"os"
	"sync"
	"woole/internal/pkg/constants"
	"woole/internal/pkg/template"
	"woole/internal/pkg/tunnel"
	"woole/pkg/rand"
	"woole/pkg/signal"
	"woole/web"

	iurl "woole/internal/pkg/url"

	"github.com/ecromaneli-golang/console/logger"
	"google.golang.org/grpc/credentials"
)

// Config has all the configuration parsed from the command line.
type Config struct {
	ClientId        string
	ClientKey       []byte
	ProxyUrl        *url.URL
	HttpUrl         *url.URL
	TunnelUrl       *url.URL
	CustomUrl       *url.URL
	SnifferUrl      *url.URL
	MaxRecords      int
	tlsSkipVerify   bool
	tlsCa           string
	EnableTLSTunnel bool
	AllowReaders    bool
	IsStandalone    bool
	available       bool
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
	proxyUrl := flag.String("proxy", ":"+defaultProxyPort, "URL of the target server to be proxied")
	tunnelUrl := flag.String("tunnel", ":"+constants.DefaultTunnelPortStr, "URL of the tunnel")
	customUrl := flag.String("custom-host", defaultCustomUrlMessage, "Custom host to be used when proxying")
	snifferPort := flag.String("sniffer", ":"+defaultSnifferPort, "Port on which the sniffing tool is available")
	maxRecords := flag.Int("records", 1000, "Max records to store. Use 0 for unlimited")
	logLevel := flag.String("log-level", "INFO", "Level of detail for the logs to be displayed")
	allowReaders := flag.Bool("allow-readers", false, "Allow other connections to listen the requests")
	tlsSkipVerify := flag.Bool("tls-skip-verify", false, "Disables the validation of the integrity of the Server's certificate")
	tlsCa := flag.String("tls-ca", "", "Path to the TLS CA file. Only for self-signed certificates")

	flag.Parse()

	logger.SetLogLevelStr(*logLevel)

	if *customUrl == defaultCustomUrlMessage {
		customUrl = proxyUrl
	}

	if *httpUrl == defaultStandaloneMessage {
		httpUrl = &emptyStr
	}

	config = &Config{
		ClientId:        *clientId,
		ClientKey:       rand.RandMD5(*clientId),
		HttpUrl:         iurl.RawUrlToUrl(*httpUrl, "http", defaultStandalonePort),
		ProxyUrl:        iurl.RawUrlToUrl(*proxyUrl, "http", ""),
		TunnelUrl:       iurl.RawUrlToUrl(*tunnelUrl, "grpc", constants.DefaultTunnelPortStr),
		CustomUrl:       iurl.RawUrlToUrl(*customUrl, "http", ""),
		SnifferUrl:      iurl.RawUrlToUrl(*snifferPort, "http", defaultSnifferPort),
		MaxRecords:      *maxRecords,
		tlsSkipVerify:   *tlsSkipVerify,
		tlsCa:           *tlsCa,
		EnableTLSTunnel: true,
		AllowReaders:    *allowReaders,
		IsStandalone:    httpUrl != &emptyStr,
		available:       true,
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
	return &tunnel.Handshake{
		ClientId:     cfg.ClientId,
		ClientKey:    cfg.ClientKey,
		AllowReaders: cfg.AllowReaders,
		Bearer:       session.Bearer,
	}
}

func PrintConfig() {
	fmt.Println(ReadConfig())
}
