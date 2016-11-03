package main

import (
	"github.com/zubairhamed/canopus"
	"log"
)

func main() {
	go runClient()
	go runServer()

	<-make(chan struct{})
}

func runClient() {
	client := canopus.NewClient()
	conn, err := client.Dial("localhost:5683")
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

	log.Println("Got Response:" + resp.GetMessage().Payload.String())
}

func runServer() {
	server := canopus.NewServer()

	server.Get("/hello", func(req canopus.CoapRequest) canopus.CoapResponse {
		log.Println("Hello Called")
		msg := canopus.ContentMessage(req.GetMessage().MessageID, canopus.MessageAcknowledgment)
		msg.SetStringPayload("Acknowledged: " + req.GetMessage().Payload.String())
		res := canopus.NewResponse(msg, nil)

		return res
	})

	server.Post("/hello", func(req canopus.CoapRequest) canopus.CoapResponse {
		log.Println("Hello Called via POST")
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

	server.ListenAndServe(":5683", nil)
}
