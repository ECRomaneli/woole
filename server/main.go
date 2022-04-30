package main

import (
	"fmt"

	"woole-server/console"
	"woole-server/recorder"
)

var config = console.ReadConfig()

func main() {
	bootstrap()
}

func bootstrap() {

	__SEPARATOR__ := separator()

	fmt.Println()
	fmt.Println(__SEPARATOR__)
	fmt.Printf("  HTTP listening on http://%s%s\n", config.HostPattern, config.HttpPort)

	tunnelProto := "http"
	if config.HasTlsFiles() {
		tunnelProto = "https"
		fmt.Printf(" HTTPS listening on https://%s%s\n", config.HostPattern, config.HttpsPort)
	}
	//fmt.Printf("Server dashboard on %s\n", config.DashboardPort)
	fmt.Printf("Tunnel listening on %s://%s%s\n", tunnelProto, config.HostPattern, config.TunnelPort)
	fmt.Println(__SEPARATOR__)
	fmt.Println()

	recorder.ListenAndServe()
}

func separator() string {
	hostLength := len(config.HostPattern) + len(config.HttpPort) + 1

	separator := "==========================="

	for i := 0; i < hostLength; i++ {
		separator += "="
	}

	return separator
}
