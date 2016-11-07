package main

import (
	"io/ioutil"
	"log"

	"github.com/zubairhamed/canopus"
)

func main() {
	go runClient()
	go runServer()

	<-make(chan struct{})
}

func runClient() {
	client := canopus.NewClient()
	conn, err := client.Dial("localhost:5683")

	file, err := ioutil.ReadFile("./ietf-block.htm")
	if err != nil {
		log.Fatal(err)
	}

	req := canopus.NewRequest(canopus.MessageConfirmable, canopus.Post, canopus.GenerateMessageID())
	blockOpt := canopus.NewBlock1Option(canopus.BlockSize16, true, 0)

	req.GetMessage().SetBlock1Option(blockOpt)
	req.SetRequestURI("/blockupload")
	req.SetPayload(file)

	resp, err := conn.Send(req)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Got Response:")
		log.Println(resp.GetMessage().GetPayload().String())
	}
}

func runServer() {
	server := canopus.NewServer()

	server.Get("/blockinfo", func(req canopus.CoapRequest) canopus.Response {
		log.Println("Hello Called")
		// canopus.PrintMessage(req.GetMessage())
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

	server.OnBlockMessage(func(msg *canopus.Message, inbound bool) {
		// log.Println("Incoming Block Message:")
		// canopus.PrintMessage(msg)
	})

	server.ListenAndServe(":5683", nil)
}
