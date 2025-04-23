package main

import (
	"fmt"

	"woole/internal/app/server/app"
	"woole/internal/app/server/recorder"
	"woole/pkg/draw"
)

var config = app.ReadConfig()

func main() {
	tunnelHost := config.GetDomain()
	if tunnelHost == "" {
		tunnelHost = "localhost"
	}

	httpUrl := fmt.Sprintf("http://%s:%s", config.HostnamePattern, config.HttpPort)
	data := []draw.KeyValue{{Key: "HTTP listening on", Value: httpUrl}}

	if config.HasTlsFiles() {
		httpsUrl := fmt.Sprintf("https://%s:%s", config.HostnamePattern, config.HttpsPort)
		data = append(data, draw.KeyValue{Key: "HTTPS listening on", Value: httpsUrl})
	}

	tunnelUrl := fmt.Sprintf("grpc://%s:%s", tunnelHost, config.TunnelPort)
	data = append(data, draw.KeyValue{Key: "Tunnel listening on", Value: tunnelUrl})

	fmt.Println("\n" + draw.Box(data))

	recorder.Start()
}
