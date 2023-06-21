package app

import (
	"crypto/tls"
	"crypto/x509"
	"embed"
	"flag"
	"fmt"
	"net/url"
	"os"
	"sync"
	"woole/internal/pkg/constants"
	pb "woole/internal/pkg/payload"
	"woole/internal/pkg/template"
	"woole/pkg/rand"
	"woole/pkg/signal"

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
	DashboardUrl    *url.URL
	MaxRecords      int
	tlsSkipVerify   bool
	tlsCa           string
	EnableTLSTunnel bool
	AllowReaders    bool
	IsStandalone    bool
	isRead          bool
}

const (
	defaultProxyPort         = "80"
	defaultDashboardPort     = "8000"
	defaultStandalonePort    = "8080"
	defaultStandaloneMessage = "[<hostname>]:<port>"
	defaultCustomUrlMessage  = "[<scheme>://]<hostname>[:<port>]"
)

//go:embed static
var EmbeddedFS embed.FS

var (
	RedirectTemplate               = template.FromFile(EmbeddedFS, "static/redirect.html")
	config           *Config       = &Config{isRead: false}
	session          *pb.Session   = &pb.Session{}
	sessionInitiated signal.Signal = *signal.New()
	writingConfig    sync.Mutex
)

func HasSession() bool {
	return session.Bearer != nil
}

// If no session was provided yet, the routine will wait for a session
func GetSessionWhenAvailable() *pb.Session {
	<-sessionInitiated.Receive()
	return session
}

func SetSession(serverSession *pb.Session) {
	if !HasSession() {
		defer sessionInitiated.SendLast()
	}
	session = serverSession
}

func ReadConfig() *Config {
	if !config.isRead {
		writingConfig.Lock()
		defer writingConfig.Unlock()
	}

	if config.isRead {
		return config
	}

	emptyStr := ""

	clientId := flag.String("client", "", "Client is an unique key used to identify the client on server")
	httpUrl := flag.String("http", defaultStandaloneMessage, "Standalone HTTP URL")
	proxyUrl := flag.String("proxy", ":"+defaultProxyPort, "URL to Proxy")
	tunnelUrl := flag.String("tunnel", ":"+constants.DefaultTunnelPortStr, "Server Tunnel URL")
	customUrl := flag.String("custom-host", defaultCustomUrlMessage, "Provide a customized URL when proxying URL")
	dashboardPort := flag.String("dashboard", ":"+defaultDashboardPort, "Dashboard Port")
	maxRecords := flag.Int("records", 16, "Max Requests to Record")
	logLevel := flag.String("log-level", "OFF", "Log Level")
	allowReaders := flag.Bool("allow-readers", false, "Allow other connections to listen the requests")
	tlsSkipVerify := flag.Bool("tls-skip-verify", false, "Do not validate the integrity of the Server's certificate")
	tlsCa := flag.String("tls-ca", "", "TLS CA file path. Only for self-signed certificates")

	flag.Parse()

	logger.SetLogLevelStr(*logLevel)

	// Deprecated: To be removed in the first stable release
	// http.DefaultTransport.(*http.Transport).TLSClientConfig = tlsConfig

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
		DashboardUrl:    iurl.RawUrlToUrl(*dashboardPort, "http", defaultDashboardPort),
		MaxRecords:      *maxRecords,
		tlsSkipVerify:   *tlsSkipVerify,
		tlsCa:           *tlsCa,
		EnableTLSTunnel: true,
		AllowReaders:    *allowReaders,
		IsStandalone:    httpUrl != &emptyStr,
		isRead:          true,
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

func (cfg *Config) GetHandshake() *pb.Handshake {
	return &pb.Handshake{
		ClientId:     cfg.ClientId,
		ClientKey:    cfg.ClientKey,
		AllowReaders: cfg.AllowReaders,
		Bearer:       session.Bearer,
	}
}

func PrintConfig() {
	fmt.Println(ReadConfig())
}
