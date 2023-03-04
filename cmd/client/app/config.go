package app

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"woole/shared/constants"
	"woole/shared/payload"
	"woole/shared/util"
	"woole/shared/util/signal"

	"github.com/ecromaneli-golang/console/logger"
)

// Config has all the configuration parsed from the command line.
type Config struct {
	ClientId     string
	ProxyUrl     *url.URL
	TunnelUrl    *url.URL
	CustomUrl    *url.URL
	DashboardUrl *url.URL
	MaxRecords   int
	isRead       bool
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
	insecureTLS := flag.Bool("allow-insecure-tls", false, "Insecure TLS verification")
	logLevel := flag.String("log-level", "OFF", "Log Level")

	flag.Parse()

	logger.SetLogLevelStr(*logLevel)

	// Ignore globally "Insecure TLS Verification"
	if *insecureTLS {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	if *customUrl == defaultCustomUrlMessage {
		customUrl = proxyUrl
	}

	config = &Config{
		ClientId:     *clientId,
		ProxyUrl:     util.RawUrlToUrl(*proxyUrl, "http", ""),
		TunnelUrl:    util.RawUrlToUrl(*tunnelUrl, "grpc", constants.DefaultTunnelPortStr),
		CustomUrl:    util.RawUrlToUrl(*customUrl, "http", ""),
		DashboardUrl: util.RawUrlToUrl(*dashboardPort, "http", defaultDashboardPort),
		MaxRecords:   *maxRecords,
		isRead:       true,
	}

	auth.ClientId = *clientId
	return config
}

func PrintConfig() {
	fmt.Println(ReadConfig())
}
