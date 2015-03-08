package main

import (
	"github.com/zubairhamed/goap"
)

/*
	Simple example to test against real CC2530/CC2538 motes over 6LowPan.
*/
func main() {
	server := goap.NewServer("udp", goap.COAP_DEFAULT_HOST)

	server.NewRoute("echo", doEcho, goap.POST)

	server.Start()
}

func doEcho(msg *goap.Message) *goap.Message {
	msg.MessageType = goap.TYPE_ACKNOWLEDGEMENT

	return msg
}
