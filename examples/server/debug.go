package main

import (
	"github.com/zubairhamed/goap"
	"log"
)

func main() {
	server := goap.NewLocalServer()

	server.NewRoute("debug", goap.GET, func(msg *goap.Message) *goap.Message {
		// goap.PrintMessage(msg)

		fwOpt := msg.GetOption(goap.OPTION_PROXY_URI)
		log.Println(fwOpt)

		ack := goap.NewMessageOfType(goap.TYPE_ACKNOWLEDGEMENT, msg.MessageId)

		return ack
	})

	server.Start()
}
