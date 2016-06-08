package main

import (
	"github.com/zubairhamed/canopus"
	"log"
)

func main() {
	server := canopus.NewLocalServer()

	server.Get("/hello", func(req canopus.CoapRequest) canopus.CoapResponse {
		log.Println("Hello Called")
		canopus.PrintMessage(req.GetMessage())
		msg := canopus.ContentMessage(req.GetMessage().MessageID, canopus.MessageAcknowledgment)
		msg.SetStringPayload("Acknowledged: " + req.GetMessage().Payload.String())
		res := canopus.NewResponse(msg, nil)

		return res
	})

	server.Post("/hello", func(req canopus.CoapRequest) canopus.CoapResponse {
		log.Println("Hello Called via POST")
		canopus.PrintMessage(req.GetMessage())
		msg := canopus.ContentMessage(req.GetMessage().MessageID, canopus.MessageAcknowledgment)
		msg.SetStringPayload("Acknowledged: " + req.GetMessage().Payload.String())
		res := canopus.NewResponse(msg, nil)

		return res
	})

	server.Get("/basic", func(req canopus.CoapRequest) canopus.CoapResponse {
		msg := canopus.NewMessageOfType(canopus.MessageAcknowledgment, req.GetMessage().MessageID)
		msg.SetStringPayload("Acknowledged")

		res := canopus.NewResponse(msg, nil)

		return res
	})

	server.Get("/basic/json", func(req canopus.CoapRequest) canopus.CoapResponse {
		msg := canopus.NewMessageOfType(canopus.MessageAcknowledgment, req.GetMessage().MessageID)
		res := canopus.NewResponse(msg, nil)

		return res
	})

	server.Get("/basic/xml", func(req canopus.CoapRequest) canopus.CoapResponse {
		msg := canopus.NewMessageOfType(canopus.MessageAcknowledgment, req.GetMessage().MessageID)
		res := canopus.NewResponse(msg, nil)

		return res
	})

	server.OnMessage(func(msg *canopus.Message, inbound bool) {
		canopus.PrintMessage(msg)
	})

	server.Start()
}
