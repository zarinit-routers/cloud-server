package main

import (
	"sync"

	"github.com/charmbracelet/log"
	"github.com/zarinit-routers/cloud-server/queue"
	"github.com/zarinit-routers/cloud-server/server"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatal(msg, "error", err)
	}
}

func main() {

	failOnError(queue.Setup(), "Failed to setup RabbitMQ connection")

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		srv := server.New()
		srv.Start()
	}()
	wg.Wait()
}
