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
	<-app.Authenticated.Receive()

	fmt.Println()
	fmt.Println("===============")
	fmt.Printf(" HTTP URL: %s\n", app.Auth.Http)

	if len(app.Auth.Https) != 0 {
		fmt.Printf("HTTPS URL: %s\n", app.Auth.Https)
	}

	fmt.Printf(" Proxying: %s\n", config.ProxyURL())
	fmt.Printf("Dashboard: http://localhost:%s\n", config.DashboardPort)

	if logger.GetInstance().IsDebugEnabled() {
		fmt.Printf("   Bearer: %s\n", app.Auth.Bearer)
	}
	fmt.Println("===============")
	fmt.Println()
}
