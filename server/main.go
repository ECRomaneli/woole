package main

import (
	"fmt"

	"woole-server/app"
	"woole-server/recorder"
)

var config = app.ReadConfig()

func main() {
	bootstrap()
}

func bootstrap() {
	fmt.Println()
	fmt.Println("===============")
	fmt.Printf("  HTTP listening on http://%s%s\n", config.HostPattern, config.HttpPort)

	tunnelProto := "http"
	if config.HasTlsFiles() {
		tunnelProto = "https"
		fmt.Printf(" HTTPS listening on https://%s%s\n", config.HostPattern, config.HttpsPort)
	}
	//fmt.Printf("Server dashboard on %s\n", config.DashboardPort)
	fmt.Printf("Tunnel listening on %s://%s%s\n", tunnelProto, config.HostPattern, config.TunnelPort)
	fmt.Println("===============")
	fmt.Println()

	recorder.ListenAndServe()
}
