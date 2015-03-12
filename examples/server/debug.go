package main

import (
	"github.com/zubairhamed/goap"
)

func main() {
	server := goap.NewLocalServer()

	server.NewRoute("debug", goap.GET, func(msg *goap.Message) *goap.Message {
		goap.PrintOptions(msg)

		ack := goap.NewMessageOfType(goap.TYPE_ACKNOWLEDGEMENT, msg.MessageId)

		return ack
	})

	server.Start()
}
