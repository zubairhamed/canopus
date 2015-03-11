package main

import (
	. "github.com/zubairhamed/goap"
)

/*
	Simple example to test against real CC2530/CC2538 motes over 6LowPan.
*/
func main() {
	server := NewLocalServer()

	server.NewRoute("event", addEvent, POST)

	server.Start()
}

func addEvent(msg *Message) *Message {
	return nil
}
