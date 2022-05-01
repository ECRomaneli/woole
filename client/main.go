package main

import (
	"fmt"
	"woole/console"
	"woole/dashboard"
	"woole/recorder"
)

func main() {
	bootstrap()
}

func bootstrap() {
	cfg := console.ReadConfig()

	go recorder.Start()

	fmt.Println()
	fmt.Println("===========================================")
	fmt.Printf("Redirecting requests from %s\n", cfg.ProxyURL())
	fmt.Printf("Serving dashboard on: http://localhost:%s\n", cfg.DashboardPort)
	fmt.Println("===========================================")
	fmt.Println()

	panic(dashboard.ListenAndServe())
}
