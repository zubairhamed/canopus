package main

import (
	. "github.com/zubairhamed/canopus"
	"log"
)

/*	To test this example, also run examples/test_server.go */
func main() {
	client := NewCoapServer(":0")

	client.OnStart(func (server *CoapServer){
		client.Dial("localhost:5683")

		req := NewRequest(TYPE_CONFIRMABLE, GET, 50782)
		req.SetStringPayload("Hello, canopus")
		req.SetRequestURI("/hello")

		resp, err := client.Send(req)
		if err != nil {
			log.Println(err)
		} else {
			log.Println("Got Synchronous Response:")
			log.Println(CoapCodeToString(resp.GetMessage().Code))
		}
	})

	client.OnMessage(func(msg *Message, inbound bool){
		if inbound {
			log.Println(">>>>> INBOUND <<<<<")
		} else {
			log.Println(">>>>> OUTBOUND <<<<<")
		}

		PrintMessage(msg)
	})

	client.Start()
}
