package main

import (
	"fmt"
	"woole/app"
	"woole/dashboard"
	"woole/recorder"
)

func main() {
	bootstrap()
}

func bootstrap() {
	config := app.ReadConfig()

	go recorder.Start()

	fmt.Println()
	fmt.Println("===========================================")
	fmt.Printf("Redirecting requests from %s\n", config.ProxyURL())
	fmt.Printf("Serving dashboard on: http://localhost:%s\n", config.DashboardPort)
	fmt.Println("===========================================")
	fmt.Println()

	panic(dashboard.ListenAndServe())
}
