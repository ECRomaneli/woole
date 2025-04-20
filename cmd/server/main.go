package main

import (
	"fmt"

	"woole/internal/app/server/app"
	"woole/internal/app/server/recorder"
)

var config = app.ReadConfig()

func main() {
	fmt.Println()
	fmt.Println("===============")
	fmt.Printf("  HTTP listening on http://%s:%s\n", config.HostnamePattern, config.HttpPort)

	if config.HasTlsFiles() {
		fmt.Printf(" HTTPS listening on https://%s:%s\n", config.HostnamePattern, config.HttpsPort)
	}

	tunnelHost := config.GetDomain()
	if tunnelHost == "" {
		tunnelHost = "localhost"
	}

	fmt.Printf("Tunnel listening on grpc://%s:%s\n", tunnelHost, config.TunnelPort)
	fmt.Println("===============")
	fmt.Println()

	recorder.Start()
}
