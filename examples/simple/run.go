package main

import (
	"log"

	"github.com/zubairhamed/canopus"
)

func main() {
	go runClient()
	go runServer()

	<-make(chan struct{})
}

func runClient() {
	conn, err := canopus.NewClient().Dial("localhost:5683")
	if err != nil {
		panic(err.Error())
	}

	req := canopus.NewRequest(canopus.MessageConfirmable, canopus.Get, canopus.GenerateMessageID())
	req.SetStringPayload("Hello, canopus")
	req.SetRequestURI("/hello")

	resp, err := conn.Send(req)
	if err != nil {
		panic(err.Error())
	}

	log.Println("Got Response:" + resp.GetMessage().GetPayload().String())
}

func runServer() {
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
		msg := canopus.NewMessageOfType(canopus.MessageAcknowledgment, req.GetMessage().GetMessageId())
		msg.SetStringPayload("Acknowledged")

		res := canopus.NewResponse(msg, nil)

		return res
	})

	server.Get("/basic/json", func(req canopus.Request) canopus.Response {
		msg := canopus.NewMessageOfType(canopus.MessageAcknowledgment, req.GetMessage().GetMessageId())
		res := canopus.NewResponse(msg, nil)

		return res
	})

	server.Get("/basic/xml", func(req canopus.Request) canopus.Response {
		msg := canopus.NewMessageOfType(canopus.MessageAcknowledgment, req.GetMessage().GetMessageId())
		res := canopus.NewResponse(msg, nil)

		return res
	})

	server.OnMessage(func(msg canopus.Message, inbound bool) {
		canopus.PrintMessage(msg)
	})

	server.ListenAndServe(":5683", nil)
}
