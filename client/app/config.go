package app

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"strings"
	"woole/util/signal"

	"github.com/ecromaneli-golang/console/logger"
)

// Config has all the configuration parsed from the command line.
type Config struct {
	ProxyProto    string
	ProxyHost     string
	ProxyPort     string
	TunnelProto   string
	TunnelHost    string
	TunnelPort    string
	DashboardPort string
	CustomHost    string
	MaxRecords    int
	isRead        bool
}

type AuthPayload struct {
	Name   string `json:"name"`
	Http   string `json:"http"`
	Https  string `json:"https"`
	Bearer string `json:"bearer"`
}

const (
	defaultProxyPort         = "80"
	defaultTunnelPort        = "8001"
	defaultDashboardPort     = "8000"
	defaultCustomHostMessage = "<scheme>://<url>:<port>"
)

func (this *Config) ProxyProtoHost() string {
	return this.ProxyProto + "://" + this.ProxyHost
}

func (this *Config) ProxyURL() string {
	if len(this.ProxyPort) == 0 {
		return this.ProxyProtoHost()
	}

	return this.ProxyProtoHost() + ":" + this.ProxyPort
}

func (this *Config) TunnelProtoHost() string {
	return this.TunnelProto + "://" + this.TunnelHost
}

func (this *Config) TunnelURL() string {
	if len(this.TunnelPort) == 0 {
		return this.TunnelProtoHost()
	}

	return this.TunnelProtoHost() + ":" + this.TunnelPort
}

var (
	config        Config        = Config{isRead: false}
	Auth          AuthPayload   = AuthPayload{}
	Authenticated signal.Signal = *signal.New()
)

// ReadConfig reads the arguments from the command line.
func ReadConfig() Config {
	if config.isRead {
		return config
	}

	proxyURL := flag.String("proxy", ":"+defaultProxyPort, "URL to Proxy")
	tunnelURL := flag.String("tunnel", ":"+defaultTunnelPort, "Server Tunnel URL. TODO: If no one is set, the sniffer will run locally")
	dashboardPort := flag.String("dashboard", defaultDashboardPort, "Dashboard Port")
	name := flag.String("name", "", "Name is an unique key used to identify the client on server")
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

	config = Config{
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

	Auth.Name = *name

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

func GetRequestURL() string {
	return fmt.Sprintf("%s/request/%s", config.TunnelURL(), Auth.Name)
}

func GetResponseURL(recordId any) string {
	return fmt.Sprintf("%s/response/%s/%s", config.TunnelURL(), Auth.Name, recordId)
}

func SetAuthorization(header http.Header) {
	header.Set("Authorization", "Bearer "+string(Auth.Bearer))
}

func PrintConfig() {
	fmt.Println(ReadConfig())
}
