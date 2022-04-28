package console

import (
	"flag"
	"fmt"

	"github.com/ecromaneli-golang/console/logger"
)

// Config has all the configuration parsed from the command line.
type Config struct {
	ServerProto    string
	ServerHost     string
	ServerPort     string
	TunnelProto    string
	TunnelHost     string
	TunnelPort     string
	DashboardProto string
	DashboardHost  string
	DashboardPort  string
	MaxRecords     int
	Timeout        int
	isRead         bool
}

const (
	defaultServerURL    = "80"
	defaultTunnelURL    = "8001"
	defaultDashboardURL = "8000"
)

var config Config = Config{isRead: false}

// ReadConfig reads the arguments from the command line.
func ReadConfig() Config {
	if config.isRead {
		return config
	}

	serverPort := flag.String("server", ":"+defaultServerURL, "Server Port")
	tunnelPort := flag.String("tunnel", ":"+defaultTunnelURL, "Tunnel Port")
	dashboardPort := flag.String("dashboard", ":"+defaultDashboardURL, "Dashboard Port")
	maxRecords := flag.Int("max", 16, "Max Requests to Record")
	timeout := flag.Int("timeout", 10000, "Timeout for receive a response from Client")
	logLevel := flag.String("log-level", "OFF", "Log Level")

	flag.Parse()

	logger.SetLogLevelStr(*logLevel)

	config = Config{
		ServerPort:    ":" + *serverPort,
		TunnelPort:    ":" + *tunnelPort,
		DashboardPort: ":" + *dashboardPort,
		MaxRecords:    *maxRecords,
		Timeout:       *timeout,
		isRead:        true,
	}

	return config
}

func PrintConfig() {
	fmt.Println(ReadConfig())
}
