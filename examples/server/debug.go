package main

import (
	. "github.com/zubairhamed/goap"
	"log"
)

func main() {
	server := NewLocalServer()

	server.NewRoute("debug", GET, func(req *CoapRequest) *CoapResponse {
		// goap.PrintMessage(msg)

		msg := req.GetMessage()
		fwOpt := msg.GetOption(OPTION_PROXY_URI)
		log.Println(fwOpt)

		ack := NewMessageOfType(TYPE_ACKNOWLEDGEMENT, msg.MessageId)

		resp := NewResponse(ack, nil)

		return resp
	})

	server.Start()
}
