package main

import (
	"github.com/zubairhamed/canopus"
	"log"
)

func main() {
	client := canopus.NewCoapServer("0")

	client.OnStart(func(server canopus.CoapServer) {
		client.Dial("localhost:5683")

		req := canopus.NewRequest(canopus.MessageConfirmable, canopus.Get, canopus.GenerateMessageID())
		req.SetStringPayload("Hello, canopus")
		req.SetRequestURI("/hello")

		for {
			resp, err := client.Send(req)
			if err != nil {
				log.Println(err)
			} else {
				log.Println("Got Response:")
				log.Println(resp.GetMessage().Payload.String())
			}
		}
	})

	client.OnError(func(err error) {
		log.Println("An error occured")
		log.Println(err)
	})

	client.OnMessage(func(msg *canopus.Message, inbound bool) {
		//if inbound {
		//	log.Println(">>>>> INBOUND <<<<<")
		//} else {
		//	log.Println(">>>>> OUTBOUND <<<<<")
		//}
		//
		//canopus.PrintMessage(msg)
	})

	client.Start()
}
