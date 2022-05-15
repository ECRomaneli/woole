package main

import (
	"fmt"
	"woole/cmd/client/app"
	"woole/cmd/client/dashboard"
	"woole/cmd/client/recorder"

	"github.com/ecromaneli-golang/console/logger"
)

var config = app.ReadConfig()

func main() {
	bootstrap()
}

func bootstrap() {
	go printInfo()
	go func() { panic(dashboard.ListenAndServe()) }()
	recorder.Start()
}

func printInfo() {
	auth := app.GetAuth()

	fmt.Println()
	fmt.Println("===============")
	fmt.Printf(" HTTP URL: %s\n", auth.HTTPUrl())

	if auth.HttpsPort != "" {
		fmt.Printf("HTTPS URL: %s\n", auth.HTTPSUrl())
	}

	fmt.Printf(" Proxying: %s\n", config.ProxyURL())
	fmt.Printf("Dashboard: %s\n", config.DashboardURL())

	if logger.GetInstance().IsDebugEnabled() {
		fmt.Printf("   Bearer: %s\n", auth.Bearer)
	}
	fmt.Println("===============")
	fmt.Println()
}
