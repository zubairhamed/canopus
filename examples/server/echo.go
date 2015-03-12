package main

import (
	. "github.com/zubairhamed/goap"
)

/*
	Simple example to test against real CC2530/CC2538 motes over 6LowPan.
*/
func main() {
	server := NewLocalServer()

	server.NewRoute("echo", doEcho, POST)

	server.Start()
}

func doEcho(msg *Message) *Message {
	msg.MessageType = TYPE_ACKNOWLEDGEMENT

	return msg
}
