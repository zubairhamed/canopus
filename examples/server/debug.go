package main

import (
	"github.com/zubairhamed/goap"
)

func main() {
	server := goap.NewLocalServer()

	server.NewRoute("debug", goap.GET, func(msg *goap.Message) *goap.Message {
		goap.PrintOptions(msg)

		msg.MessageType = goap.TYPE_ACKNOWLEDGEMENT
		return msg
	})

	server.Start()
}
