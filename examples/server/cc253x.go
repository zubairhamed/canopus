package main

import (
	"github.com/zubairhamed/goap"
)

/*
	Simple example to test against real CC2530/CC2538 motes over 6LowPan.
*/
func main() {
	server := goap.NewLocalServer()

	server.NewRoute("event", addEvent, goap.POST)

	server.Start()
}

func addEvent(msg *goap.Message) *goap.Message {
	return nil
}
