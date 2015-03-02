package main

import (
	"github.com/zubairhamed/goap"
)

/*
	Simple example to test against real CC2530/CC2538 motes over 6LowPan.
 */
func main() {
	server := goap.NewServer("udp", goap.COAP_DEFAULT_HOST)

	server.NewRoute("state", addState, goap.POST)

	server.Start()
}

func addState(msg *goap.Message) *goap.Message {
	return nil
}
