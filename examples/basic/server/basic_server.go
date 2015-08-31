package main

import (
	. "github.com/zubairhamed/canopus"
	"log"
)

func main() {
	server := NewLocalServer()

	server.Get("/hello", func(req CoapRequest) CoapResponse {
		log.Println("Hello Called")
		PrintMessage(req.GetMessage())
		msg := ContentMessage(req.GetMessage().MessageId, TYPE_ACKNOWLEDGEMENT)
		msg.SetStringPayload("Acknowledged: " + req.GetMessage().Payload.String())
		res := NewResponse(msg, nil)

		return res
	})

	server.Post("/hello", func(req CoapRequest) CoapResponse {
		log.Println("Hello Called via POST")
		PrintMessage(req.GetMessage())
		msg := ContentMessage(req.GetMessage().MessageId, TYPE_ACKNOWLEDGEMENT)
		msg.SetStringPayload("Acknowledged: " + req.GetMessage().Payload.String())
		res := NewResponse(msg, nil)

		return res
	})


	server.Get("/basic", func(req CoapRequest) CoapResponse {
		msg := NewMessageOfType(TYPE_ACKNOWLEDGEMENT, req.GetMessage().MessageId)
		msg.SetStringPayload("Acknowledged")

		res := NewResponse(msg, nil)

		return res
	})

	server.Get("/basic/json", func(req CoapRequest) CoapResponse {
		msg := NewMessageOfType(TYPE_ACKNOWLEDGEMENT, req.GetMessage().MessageId)
		res := NewResponse(msg, nil)

		return res
	})

	server.Get("/basic/xml", func(req CoapRequest) CoapResponse {
		msg := NewMessageOfType(TYPE_ACKNOWLEDGEMENT, req.GetMessage().MessageId)
		res := NewResponse(msg, nil)

		return res
	})

	server.OnMessage(func(msg *Message, inbound bool){
		PrintMessage(msg)
	})

	server.Start()
}
