package console

import (
	"flag"
	"fmt"

	"github.com/ecromaneli-golang/console/logger"
)

// Config has all the configuration parsed from the command line.
type Config struct {
	ServerPort    string
	TunnelPort    string
	DashboardPort string
	Timeout       int
	TlsFullChain  string
	TlsPrivKey    string
	isRead        bool
}

const (
	defaultServerPort    = "80"
	defaultTunnelPort    = "8001"
	defaultDashboardPort = "8000"
)

var config Config = Config{isRead: false}

// ReadConfig reads the arguments from the command line.
func ReadConfig() Config {
	if config.isRead {
		return config
	}

	serverPort := flag.String("server", defaultServerPort, "Server Port")
	tunnelPort := flag.String("tunnel", defaultTunnelPort, "Tunnel Port")
	dashboardPort := flag.String("dashboard", defaultDashboardPort, "Dashboard Port")
	timeout := flag.Int("timeout", 10000, "Timeout for receive a response from Client")
	tlsFullChain := flag.String("tls-fullchain", "", "TLS fullchain file path")
	tlsPrivKey := flag.String("tls-privkey", "", "TLS privkey file path")
	logLevel := flag.String("log-level", "OFF", "Log Level")

	flag.Parse()

	logger.SetLogLevelStr(*logLevel)

	config = Config{
		ServerPort:    ":" + *serverPort,
		TunnelPort:    ":" + *tunnelPort,
		DashboardPort: ":" + *dashboardPort,
		TlsFullChain:  *tlsFullChain,
		TlsPrivKey:    *tlsPrivKey,
		Timeout:       *timeout,
		isRead:        true,
	}

	return config
}

func PrintConfig() {
	fmt.Println(ReadConfig())
}
