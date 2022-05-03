package app

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/ecromaneli-golang/console/logger"
)

const ClientToken = "{client}"

// Config has all the configuration parsed from the command line.
type Config struct {
	HostPattern   string
	HttpPort      string
	HttpsPort     string
	TunnelPort    string
	DashboardPort string
	Timeout       int
	TlsCert       string
	TlsKey        string
	isRead        bool
}

var config Config = Config{isRead: false}

func (this *Config) HasTlsFiles() bool {
	return len(this.TlsCert) != 0 && len(this.TlsKey) != 0
}

// ReadConfig reads the arguments from the command line.
func ReadConfig() Config {
	if config.isRead {
		return config
	}

	hostPattern := flag.String("pattern", "{client}", "Set the server host pattern. The '"+ClientToken+"' MUST be present to determine where to get client id")
	httpPort := flag.Int("http", 80, "HTTP Port")
	httpsPort := flag.Int("https", 443, "HTTPS Port")
	tunnelPort := flag.Int("tunnel", 8001, "Tunnel Port")
	dashboardPort := flag.Int("dashboard", 8000, "Dashboard Port")
	timeout := flag.Int("timeout", 10000, "Timeout for receive a response from Client")
	tlsCert := flag.String("tls-cert", "", "TLS cert/fullchain file path")
	tlsKey := flag.String("tls-key", "", "TLS key/privkey file path")
	logLevel := flag.String("log-level", "OFF", "Log Level")

	flag.Parse()

	logger.SetLogLevelStr(*logLevel)

	config = Config{
		HostPattern:   *hostPattern,
		HttpPort:      ":" + strconv.Itoa(*httpPort),
		HttpsPort:     ":" + strconv.Itoa(*httpsPort),
		TunnelPort:    ":" + strconv.Itoa(*tunnelPort),
		DashboardPort: ":" + strconv.Itoa(*dashboardPort),
		TlsCert:       *tlsCert,
		TlsKey:        *tlsKey,
		Timeout:       *timeout,
		isRead:        true,
	}

	if strings.Index(config.HostPattern, ClientToken) == -1 {
		panic("Pattern MUST has " + ClientToken)
	}

	return config
}

func PrintConfig() {
	fmt.Println(ReadConfig())
}
