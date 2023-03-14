package main

import (
	"fmt"
	"woole/cmd/client/app"
	"woole/cmd/client/dashboard"
	"woole/cmd/client/recorder"
)

var config = app.ReadConfig()

func main() {
	bootstrap()
}

func bootstrap() {
	go printInfo()
	go startDashboard()
	recorder.Start()
}

func startDashboard() {
	panic(dashboard.ListenAndServe())
}

func printInfo() {
	auth := app.GetSession()

	fmt.Println()
	fmt.Println("===============")
	fmt.Printf(" HTTP URL: %s\n", auth.HttpUrl())

	if auth.HttpsPort != "" {
		fmt.Printf("HTTPS URL: %s\n", auth.HttpsUrl())
	}

	fmt.Printf(" Proxying: %s\n", config.ProxyUrl.String())
	fmt.Printf("Dashboard: %s\n", config.DashboardUrl.String())

	fmt.Println("===============")
	fmt.Println()
}
