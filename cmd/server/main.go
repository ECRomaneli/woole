package main

import (
	"fmt"

	"woole/cmd/server/app"
	"woole/cmd/server/recorder"
)

var config = app.ReadConfig()

func main() {
	bootstrap()
}

func bootstrap() {
	fmt.Println()
	fmt.Println("===============")
	fmt.Printf("  HTTP listening on http://%s:%s\n", config.HostnamePattern, config.HttpPort)

	if config.HasTlsFiles() {
		fmt.Printf(" HTTPS listening on https://%s:%s\n", config.HostnamePattern, config.HttpsPort)
	}

	fmt.Printf("Tunnel listening on grpc://%s:%s\n", config.HostnamePattern, config.TunnelPort)
	fmt.Println("===============")
	fmt.Println()

	recorder.Start()
}
