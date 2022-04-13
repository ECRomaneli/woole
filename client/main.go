package main

import (
	"fmt"
	"sync"
	"woole/console"
	"woole/dashboard"
	"woole/recorder"
)

func main() {
	bootstrap()
}

func bootstrap() {
	cfg := console.ReadConfig()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		recorder.Start()
		wg.Done()
	}()

	go func() {
		err := dashboard.ListenAndServe()
		fmt.Println("Dashboard: ", err)
		wg.Done()
	}()

	fmt.Println()
	fmt.Println("===========================================")
	fmt.Printf("Redirecting requests from %s\n", cfg.ProxyURL())
	fmt.Printf("Serving dashboard on: http://localhost:%s\n", cfg.DashboardPort)
	fmt.Println("===========================================")
	fmt.Println()

	wg.Wait()
}
