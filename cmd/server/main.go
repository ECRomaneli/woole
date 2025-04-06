package main

import (
	"fmt"
	"net"

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

	localIp, err := getLocalIp()
	if err != nil {
		fmt.Printf("Tunnel listening on grpc://<hostname-or-ip>:%s\n", config.TunnelPort)
	} else {
		fmt.Printf("Tunnel listening on grpc://%s:%s\n", localIp, config.TunnelPort)
	}
	fmt.Println("===============")
	fmt.Println()

	recorder.Start()
}

func getLocalIp() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil { // Verifica se é IPv4
				return ipNet.IP.String(), nil
			}
		}
	}

	return "", fmt.Errorf("failed to retrieve the local ip")
}
