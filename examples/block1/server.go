package main

import (
	"io/ioutil"
	"log"

	"github.com/zubairhamed/canopus"
)

func main() {
	server := canopus.NewServer()

	server.Get("/blockinfo", func(req canopus.Request) canopus.Response {
		msg := canopus.ContentMessage(req.GetMessage().GetMessageId(), canopus.MessageAcknowledgment)
		msg.SetStringPayload("Acknowledged: " + req.GetMessage().GetPayload().String())

		res := canopus.NewResponse(msg, nil)

		return res
	})

	server.Post("/blockupload", func(req canopus.Request) canopus.Response {
		msg := canopus.ContentMessage(req.GetMessage().GetMessageId(), canopus.MessageAcknowledgment)
		msg.SetStringPayload("Acknowledged: " + req.GetMessage().GetPayload().String())
		res := canopus.NewResponse(msg, nil)

		// Save to file
		payload := req.GetMessage().GetPayload().GetBytes()
		log.Println("len", len(payload))
		err := ioutil.WriteFile("output.html", payload, 0644)
		if err != nil {
			log.Println(err)
		}

		return res
	})

	server.OnBlockMessage(func(msg canopus.Message, inbound bool) {
		// log.Println("Incoming Block Message:")
		// canopus.PrintMessage(msg)
	})

	server.ListenAndServe(":5683")
	<-make(chan struct{})
}
