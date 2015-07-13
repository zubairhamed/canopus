package main

import (
	. "github.com/zubairhamed/canopus"
	"log"
)

/*	To test this example, also run examples/test_server.go */
func main() {
	log.Println("Starting Client..")
	client := NewCoapServer(":0")

	client.On(EVT_START, func(){
		log.Println("EVT_START")
		client.Dial(":5683")

		req := NewRequest(TYPE_CONFIRMABLE, GET, 50782)
		req.SetStringPayload("Hello, canopus")
		req.SetRequestURI("/0/1/2")

		resp, err := client.Send(req)
		if err != nil {
			log.Println(err)
		} else {
			log.Println("Got Synchronous Response:")
			log.Println(CoapCodeToString(resp.GetMessage().Code))
		}
	})
	client.Start()
}
