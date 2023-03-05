package app

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"net/url"
	"sync"
	"woole/shared/constants"
	"woole/shared/payload"
	"woole/shared/util"
	"woole/shared/util/signal"

	"github.com/ecromaneli-golang/console/logger"
	"google.golang.org/grpc/credentials"
)

// Config has all the configuration parsed from the command line.
type Config struct {
	ClientId      string
	ProxyUrl      *url.URL
	TunnelUrl     *url.URL
	CustomUrl     *url.URL
	DashboardUrl  *url.URL
	MaxRecords    int
	tlsSkipVerify bool
	tlsCa         string
	isRead        bool
}

const (
	defaultProxyPort        = "80"
	defaultDashboardPort    = "8000"
	defaultCustomUrlMessage = "[<scheme>://]<hostname>[:<port>]"
)

var (
	config        *Config       = &Config{isRead: false}
	auth          *payload.Auth = &payload.Auth{}
	authenticated signal.Signal = *signal.New()
	writingConfig sync.Mutex
)

// If no auth was provided yet, the routine will wait for an authentication
func GetAuth() *payload.Auth {
	<-authenticated.Receive()
	return auth
}

func Authenticate(authentication *payload.Auth) {
	auth = authentication
	if auth.Bearer == "" {
		panic("No bearer was provided")
	}
	authenticated.SendLast()
}

func ReadConfig() *Config {
	if !config.isRead {
		writingConfig.Lock()
		defer writingConfig.Unlock()
	}

	if config.isRead {
		return config
	}

	clientId := flag.String("client", "", "Client is an unique key used to identify the client on server")
	proxyUrl := flag.String("proxy", ":"+defaultProxyPort, "URL to Proxy")
	tunnelUrl := flag.String("tunnel", ":"+constants.DefaultTunnelPortStr, "Server Tunnel URL") // TODO: If no one is set, the sniffer will run locally
	customUrl := flag.String("custom-host", defaultCustomUrlMessage, "Provide a customized URL when proxying URL")
	dashboardPort := flag.String("dashboard", ":"+defaultDashboardPort, "Dashboard Port")
	maxRecords := flag.Int("records", 16, "Max Requests to Record")
	logLevel := flag.String("log-level", "OFF", "Log Level")
	tlsSkipVerify := flag.Bool("tls-skip-verify", false, "Do not validate the integrity of the Server's certificate")
	tlsCa := flag.String("tls-ca", "", "TLS CA file path. Only for self-signed certificates")

	flag.Parse()

	logger.SetLogLevelStr(*logLevel)

	// Deprecated: To be removed in the first stable release
	// http.DefaultTransport.(*http.Transport).TLSClientConfig = tlsConfig

	if *customUrl == defaultCustomUrlMessage {
		customUrl = proxyUrl
	}

	config = &Config{
		ClientId:      *clientId,
		ProxyUrl:      util.RawUrlToUrl(*proxyUrl, "http", ""),
		TunnelUrl:     util.RawUrlToUrl(*tunnelUrl, "grpc", constants.DefaultTunnelPortStr),
		CustomUrl:     util.RawUrlToUrl(*customUrl, "http", ""),
		DashboardUrl:  util.RawUrlToUrl(*dashboardPort, "http", defaultDashboardPort),
		MaxRecords:    *maxRecords,
		tlsSkipVerify: *tlsSkipVerify,
		tlsCa:         *tlsCa,
		isRead:        true,
	}

	auth.ClientId = *clientId
	return config
}

func (cfg *Config) GetTransportCredentials() credentials.TransportCredentials {
	tlsConfig := &tls.Config{InsecureSkipVerify: cfg.tlsSkipVerify}

	if tlsConfig.InsecureSkipVerify || cfg.tlsCa == "" {
		return credentials.NewTLS(tlsConfig)
	}

	caBytes, err := ioutil.ReadFile(cfg.tlsCa)
	if err != nil {
		panic("Failed to load TLS CA File on: " + cfg.tlsCa + ". Reason: " + err.Error())
	}

	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM(caBytes) {
		panic("Failed to append certificate")
	}

	return credentials.NewTLS(tlsConfig)
}

func PrintConfig() {
	fmt.Println(ReadConfig())
}
