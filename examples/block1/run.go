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
	req.SetPayload(file)
	req.SetRequestURI("/blockupload")

	resp, err := conn.Send(req)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Got Response:")
		log.Println(resp.GetMessage().Payload.String())
	}
}

func runServer() {
	server := canopus.NewLocalServer("TestServer")

	server.Get("/blockinfo", func(req canopus.CoapRequest) canopus.CoapResponse {
		log.Println("Hello Called")
		// canopus.PrintMessage(req.GetMessage())
		msg := canopus.ContentMessage(req.GetMessage().MessageID, canopus.MessageAcknowledgment)
		msg.SetStringPayload("Acknowledged: " + req.GetMessage().Payload.String())
		res := canopus.NewResponse(msg, nil)

		return res
	})

	server.Post("/blockupload", func(req canopus.CoapRequest) canopus.CoapResponse {
		msg := canopus.ContentMessage(req.GetMessage().MessageID, canopus.MessageAcknowledgment)
		msg.SetStringPayload("Acknowledged: " + req.GetMessage().Payload.String())
		res := canopus.NewResponse(msg, nil)

		// Save to file

		payload := req.GetMessage().Payload.GetBytes()
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

	server.Start()
}
