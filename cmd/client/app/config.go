package app

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"woole/shared/payload"
	"woole/shared/util/signal"

	"github.com/ecromaneli-golang/console/logger"
)

// Config has all the configuration parsed from the command line.
type Config struct {
	ClientId      string
	ProxyProto    string
	ProxyHost     string
	ProxyPort     string
	TunnelProto   string
	TunnelHost    string
	TunnelPort    string
	DashboardPort string
	CustomHost    string
	MaxRecords    int
	mu            sync.Mutex
	isRead        bool
}

const (
	defaultProxyPort         = "80"
	defaultTunnelPort        = "8001"
	defaultDashboardPort     = "8000"
	defaultCustomHostMessage = "<scheme>://<url>:<port>"
)

func (cfg *Config) ProxyProtoHost() string {
	return cfg.ProxyProto + "://" + cfg.ProxyHost
}

func (cfg *Config) ProxyURL() string {
	if len(cfg.ProxyPort) == 0 {
		return cfg.ProxyProtoHost()
	}

	return cfg.ProxyProtoHost() + ":" + cfg.ProxyPort
}

func (cfg *Config) TunnelProtoHost() string {
	return cfg.TunnelProto + "://" + cfg.TunnelHost
}

func (cfg *Config) TunnelURL() string {
	if len(cfg.TunnelPort) == 0 {
		return cfg.TunnelProtoHost()
	}

	return cfg.TunnelProtoHost() + ":" + cfg.TunnelPort
}

func (cfg *Config) TunnelHostPort() string {
	if len(cfg.TunnelPort) == 0 {
		return cfg.TunnelHost
	}

	return cfg.TunnelHost + ":" + cfg.TunnelPort
}

func (cfg *Config) DashboardURL() string {
	return "http://localhost:" + cfg.DashboardPort
}

var (
	config        *Config       = &Config{isRead: false}
	auth          *payload.Auth = &payload.Auth{}
	authenticated signal.Signal = *signal.New()
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
		config.mu.Lock()
		defer config.mu.Unlock()
	}

	if config.isRead {
		return config
	}

	proxyURL := flag.String("proxy", ":"+defaultProxyPort, "URL to Proxy")
	tunnelURL := flag.String("tunnel", ":"+defaultTunnelPort, "Server Tunnel URL. TODO: If no one is set, the sniffer will run locally")
	dashboardPort := flag.String("dashboard", defaultDashboardPort, "Dashboard Port")
	client := flag.String("client", "", "Client is an unique key used to identify the client on server")
	customHost := flag.String("custom-host", defaultCustomHostMessage, "Customize host passed as header for proxy URL")
	maxRecords := flag.Int("records", 16, "Max Requests to Record")
	insecureTLS := flag.Bool("allow-insecure-tls", false, "Insecure TLS verification")
	logLevel := flag.String("log-level", "OFF", "Log Level")

	flag.Parse()

	logger.SetLogLevelStr(*logLevel)

	// Ignore globally "Insecure TLS Verification"
	if *insecureTLS {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	proxyProto, proxyHost, proxyPort := splitURL(*proxyURL)
	tunnelProto, tunnelHost, tunnelPort := splitURL(*tunnelURL)

	config = &Config{
		ClientId:      *client,
		ProxyProto:    strOrDefault(proxyProto, "http"),
		ProxyHost:     strOrDefault(proxyHost, "localhost"),
		ProxyPort:     proxyPort,
		TunnelProto:   strOrDefault(tunnelProto, "http"),
		TunnelHost:    strOrDefault(tunnelHost, "localhost"),
		TunnelPort:    tunnelPort,
		DashboardPort: *dashboardPort,
		CustomHost:    *customHost,
		MaxRecords:    *maxRecords,
		isRead:        true,
	}

	auth.ClientId = *client

	if config.CustomHost == defaultCustomHostMessage {
		config.CustomHost = config.ProxyURL()
	}

	return config
}

func strOrDefault(val string, def string) string {
	if len(val) > 0 {
		return val
	}
	return def
}

func splitURL(url string) (proto, host, port string) {
	host = url

	colon := strings.Index(host, "://")
	if colon == -1 {
		host, port = splitHostPort(host)
		return "", host, port
	}

	host, port = splitHostPort(url[colon+3:])
	return url[:colon], host, port
}

func splitHostPort(hostPort string) (host, port string) {
	host = hostPort

	colon := strings.LastIndexByte(host, ':')
	if colon == -1 {
		return host, ""
	}

	return hostPort[:colon], hostPort[colon+1:]
}

func PrintConfig() {
	fmt.Println(ReadConfig())
}
