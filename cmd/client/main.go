package main

import (
	"fmt"
	"woole/internal/app/client/app"
	"woole/internal/app/client/recorder"
	"woole/internal/app/client/sniffer"
)

var config = app.ReadConfig()

func main() {
	go printInfo()
	go startSnifferTool()
	recorder.Start()
}

func startSnifferTool() {
	panic(sniffer.ListenAndServe())
}

func printInfo() {
	session := app.GetSessionWhenAvailable()

	fmt.Println()
	fmt.Println("===============")
	fmt.Printf(" HTTP URL: %s\n", session.HttpUrl())

	if session.HttpsPort != "" {
		fmt.Printf("HTTPS URL: %s\n", session.HttpsUrl())
	}

	fmt.Printf(" Proxying: %s\n", config.ProxyUrl.String())
	fmt.Printf("  Sniffer: %s\n", config.SnifferUrl.String())
	fmt.Printf("Expire At: %s\n", app.ExpireAt())

	fmt.Println("===============")
	fmt.Println()
}
