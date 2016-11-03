package main

import (
	"log"

	"github.com/zubairhamed/canopus"
)

func main() {
	go runClient()
	// go runServer()

	<-make(chan struct{})
}

func runClient() {
	client := canopus.NewClient()
	conn, err := client.DialDTLS("localhost:5684", "secretPSK")
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

	cfg := &canopus.ServerConfiguration{}

	server.ListenAndServeDTLS(":5684", cfg)
}
