package main

import (
	"fmt"
	"os"

	"woole/internal/app/server/app"
	"woole/internal/app/server/recorder"
	"woole/pkg/draw"

	"github.com/ecromaneli-golang/console/logger"
)

var config = app.ReadConfig()

func main() {
	defer func() {
		if err := recover(); err != nil {
			log := logger.New("woole")
			log.SetLogLevelStr(config.LogLevel)
			log.Fatal(err)
			if log.IsDebugEnabled() {
				panic(err)
			}
			os.Exit(1)
		}
	}()

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
