package main

import (
	"fmt"
	"woole/internal/app/client/app"
	"woole/internal/app/client/recorder"
	"woole/internal/app/client/sniffer"
	"woole/pkg/draw"
)

var config = app.ReadConfig()

func main() {
	go printInfo()
	go recorder.Start()
	startSnifferTool()
}

func startSnifferTool() {
	panic(sniffer.ListenAndServe())
}

func printInfo() {
	session := app.GetSessionWhenAvailable()

	data := []draw.KeyValue{
		{Key: "HTTP URL", Value: session.HttpUrl()},
	}

	if session.HttpsPort != "" {
		data = append(data, draw.KeyValue{Key: "HTTPS URL", Value: session.HttpsUrl()})
	}

	data = append(data,
		draw.KeyValue{Key: "Proxying", Value: config.ProxyUrl.String()},
		draw.KeyValue{Key: "Sniffer", Value: config.SnifferUrl.String()},
		draw.KeyValue{Key: "Expire At", Value: app.ExpireAt()})

	fmt.Println("\n" + draw.Box(data))
}
