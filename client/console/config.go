package console

import (
	"flag"
	"fmt"
	"strings"

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
	MaxRecords    int
	isRead        bool
}

const (
	defaultProxyURL     = "80"
	defaultTunnelURL    = "8001"
	defaultDashboardURL = "8000"
)

func (this *Config) ProxyProtoHost() string {
	return this.ProxyProto + "://" + this.ProxyHost
}

func (this *Config) ProxyURL() string {
	return this.ProxyProtoHost() + ":" + this.ProxyPort
}

func (this *Config) TunnelProtoHost() string {
	return this.TunnelProto + "://" + this.TunnelHost
}

func (this *Config) TunnelURL() string {
	return this.TunnelProtoHost() + ":" + this.TunnelPort
}

var config Config = Config{isRead: false}

// ReadConfig reads the arguments from the command line.
func ReadConfig() Config {
	if config.isRead {
		return config
	}

	proxyURL := flag.String("proxy", ":"+defaultProxyURL, "URL to Proxy")
	tunnelURL := flag.String("tunnel", ":"+defaultTunnelURL, "Server Tunnel URL. TODO: If no one is set, the sniffer will run locally")
	dashboardPort := flag.String("dashboard", defaultDashboardURL, "Dashboard Port")
	maxRecords := flag.Int("records", 16, "Max Requests to Record")
	logLevel := flag.String("log-level", "OFF", "Log Level")

	flag.Parse()

	logger.SetLogLevelStr(*logLevel)

	proxyProto, proxyHost, proxyPort := splitURL(*proxyURL)
	tunnelProto, tunnelHost, tunnelPort := splitURL(*tunnelURL)

	config = Config{
		ProxyProto:    strOrDefault(proxyProto, "http"),
		ProxyHost:     strOrDefault(proxyHost, "localhost"),
		ProxyPort:     strOrDefault(proxyPort, defaultProxyURL),
		TunnelProto:   strOrDefault(tunnelProto, "http"),
		TunnelHost:    strOrDefault(tunnelHost, "localhost"),
		TunnelPort:    strOrDefault(tunnelPort, defaultTunnelURL),
		DashboardPort: *dashboardPort,
		MaxRecords:    *maxRecords,
		isRead:        true,
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
