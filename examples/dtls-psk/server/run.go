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

	server.HandlePSK(func(id string) []byte {
		return []byte("secretPSK")
	})

	server.ListenAndServeDTLS(":5684")

	<-make(chan struct{})
}
