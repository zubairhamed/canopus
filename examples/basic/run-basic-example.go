package main

import (
	"github.com/zubairhamed/canopus"
	"log"
)

func main() {
	go runServer()
	runClient()

	for {

	}
}

func runServer() {
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

func runClient() {
	client := canopus.NewCoapServer("0")

	client.OnStart(func(server canopus.CoapServer) {
		client.Dial("localhost:5683")

		req := canopus.NewRequest(canopus.MessageConfirmable, canopus.Get, canopus.GenerateMessageID())
		req.GetMessage().AddOption(canopus.OptionBlock1, nil)
		req.SetStringPayload("Hello, canopus")
		req.SetRequestURI("/hello")

		for {
			_, err := client.Send(req)
			if err != nil {
			log.Println(err)
			} else {
			log.Println("Got Response:")
			// log.Println(resp.GetMessage().Payload.String())
			}
		}

	})

	client.OnError(func(err error) {
		log.Println("An error occured")
		log.Println(err)
	})

	client.OnMessage(func(msg *canopus.Message, inbound bool) {
		if inbound {
			log.Println(">>>>> INBOUND <<<<<")
		} else {
			log.Println(">>>>> OUTBOUND <<<<<")
		}

		canopus.PrintMessage(msg)
	})

	client.Start()
}
