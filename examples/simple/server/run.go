package main

import (
	"log"

	"github.com/zubairhamed/canopus"
)

func main() {
	server := canopus.NewServer()

	server.Get("/hello", func(req canopus.Request) canopus.Response {
		log.Println("Hello Called")
		msg := canopus.ContentMessage(req.GetMessage().GetMessageId(), canopus.MessageAcknowledgment)
		msg.SetStringPayload("Acknowledged: " + req.GetMessage().GetPayload().String())

		res := canopus.NewResponse(msg, nil)
		return res
	})

	server.Post("/hello", func(req canopus.Request) canopus.Response {
		log.Println("Hello Called via POST")

		msg := canopus.ContentMessage(req.GetMessage().GetMessageId(), canopus.MessageAcknowledgment)
		msg.SetStringPayload("Acknowledged: " + req.GetMessage().GetPayload().String())
		res := canopus.NewResponse(msg, nil)
		return res
	})

	server.Get("/basic", func(req canopus.Request) canopus.Response {
		msg := canopus.NewMessageOfType(canopus.MessageAcknowledgment, req.GetMessage().GetMessageId(), canopus.NewPlainTextPayload("Acknowledged"))
		res := canopus.NewResponse(msg, nil)
		return res
	})

	server.Get("/basic/json", func(req canopus.Request) canopus.Response {
		msg := canopus.NewMessageOfType(canopus.MessageAcknowledgment, req.GetMessage().GetMessageId(), nil)

		res := canopus.NewResponse(msg, nil)

		return res
	})

	server.Get("/basic/xml", func(req canopus.Request) canopus.Response {
		msg := canopus.NewMessageOfType(canopus.MessageAcknowledgment, req.GetMessage().GetMessageId(), nil)
		res := canopus.NewResponse(msg, nil)

		return res
	})

	server.OnMessage(func(msg canopus.Message, inbound bool) {
		canopus.PrintMessage(msg)
	})

	server.ListenAndServe(":5683")
	<-make(chan struct{})
}
