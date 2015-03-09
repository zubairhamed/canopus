package main

import (
	"github.com/zubairhamed/goap"
)

func main() {
	server := goap.NewLocalServer()

	server.NewRoute("serviceA", goap.GET, service)
	server.NewRoute("serviceB", goap.GET, service)
	server.NewRoute("serviceC", goap.GET, service)
	server.NewRoute("serviceD", goap.GET, service)

	server.Start()
}

func service(msg *goap.Message) *goap.Message {
	return nil
}
