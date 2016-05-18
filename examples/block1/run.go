package main

import (
	"github.com/zubairhamed/canopus"
	"io/ioutil"
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

	server.Get("/blockinfo", func(req canopus.CoapRequest) canopus.CoapResponse {
		log.Println("Hello Called")
		// canopus.PrintMessage(req.GetMessage())
		msg := canopus.ContentMessage(req.GetMessage().MessageID, canopus.MessageAcknowledgment)
		msg.SetStringPayload("Acknowledged: " + req.GetMessage().Payload.String())
		res := canopus.NewResponse(msg, nil)

		return res
	})

	server.Post("/blockupload", func(req canopus.CoapRequest) canopus.CoapResponse {
		log.Println("Hello Called via POST")
		// canopus.PrintMessage(req.GetMessage())
		msg := canopus.ContentMessage(req.GetMessage().MessageID, canopus.MessageAcknowledgment)
		msg.SetStringPayload("Acknowledged: " + req.GetMessage().Payload.String())
		res := canopus.NewResponse(msg, nil)

		return res
	})

	server.OnBlockMessage(func(msg *canopus.Message, inbound bool) {
		log.Println("Incoming Block Message:")
		canopus.PrintMessage(msg)
	})

	server.Start()
}

func runClient() {
	client := canopus.NewCoapServer("0")

	client.OnStart(func(server canopus.CoapServer) {
		client.Dial("localhost:5683")

		file, err := ioutil.ReadFile("./feedbooks_book_1421.epub")
		if err != nil {

		}

		req := canopus.NewRequest(canopus.MessageConfirmable, canopus.Get, canopus.GenerateMessageID())

		// blockOpt := canopus.NewBlock1Option(canopus.BlockSize16, true)
		blockOpt := canopus.NewBlock1Option(canopus.BlockSize64, true)

		//  req.GetMessage().AddOption(canopus.OptionBlock1, nil)
		req.GetMessage().SetBlock1Option(blockOpt)
		// req.SetStringPayload("Returning Block transfer information")
		req.SetPayload(file)
		req.SetRequestURI("/blockinfo")

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

	//client.OnMessage(func(msg *canopus.Message, inbound bool) {
	//	if inbound {
	//		log.Println(">>>>> INBOUND <<<<<")
	//	} else {
	//		log.Println(">>>>> OUTBOUND <<<<<")
	//	}
	//
	//	canopus.PrintMessage(msg)
	//})

	client.Start()
}
