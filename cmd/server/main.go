package main

import (
	"sync"

	"github.com/charmbracelet/log"
	"github.com/zarinit-routers/cloud-server/server"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatal(msg, "error", err)
	}
}

func main() {

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		srv := server.New()
		srv.Start()
	}()
	wg.Wait()
}
