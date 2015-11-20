package main

import (
	. "github.com/zubairhamed/canopus"
	"log"
)

func main() {
	go runServer()
	runClient()

	for {

	}
}

func runServer() {
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

	server.OnMessage(func(msg *Message, inbound bool) {
		PrintMessage(msg)
	})

	server.Start()
}

func runClient() {
	client := NewCoapServer(":0")

	client.OnStart(func(server CoapServer) {
		client.Dial("localhost:5683")

		req := NewRequest(TYPE_CONFIRMABLE, GET, GenerateMessageId())
		req.SetStringPayload("Hello, canopus")
		req.SetRequestURI("/hello")

		resp, err := client.Send(req)
		if err != nil {
			log.Println(err)
		} else {
			log.Println("Got Response:")
			log.Println(resp.GetMessage().Payload.String())
		}
	})

	client.OnError(func(err error) {
		log.Println("An error occured")
		log.Println(err)
	})

	client.OnMessage(func(msg *Message, inbound bool) {
		if inbound {
			log.Println(">>>>> INBOUND <<<<<")
		} else {
			log.Println(">>>>> OUTBOUND <<<<<")
		}

		PrintMessage(msg)
	})

	client.Start()
}
