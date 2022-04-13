package main

import (
	"fmt"
	"sync"

	"woole-server/console"
	"woole-server/recorder"
)

func main() {
	bootstrap()
}

func bootstrap() {
	cfg := console.ReadConfig()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		err := recorder.ListenAndServe()
		fmt.Println("Recorder: ", err)
		wg.Done()
	}()

	fmt.Println()
	fmt.Println("=========================")
	fmt.Printf("Server listening on %s\n", cfg.ServerPort)
	//fmt.Printf("Server dashboard on %s\n", cfg.DashboardPort)
	fmt.Printf("Tunnel listening on %s\n", cfg.TunnelPort)
	fmt.Println("=========================")
	fmt.Println()

	wg.Wait()
}
