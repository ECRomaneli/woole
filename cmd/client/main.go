package main

import (
	"fmt"
	"woole/internal/app/client/app"
	"woole/internal/app/client/dashboard"
	"woole/internal/app/client/recorder"
)

var config = app.ReadConfig()

func main() {
	go printInfo()
	go startDashboard()
	recorder.Start()
}

func startDashboard() {
	panic(dashboard.ListenAndServe())
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
	fmt.Printf("Dashboard: %s\n", config.DashboardUrl.String())

	fmt.Println("===============")
	fmt.Println()
}
