package main

import (
	. "github.com/zubairhamed/canopus"
	"log"
)

func main() {
	client := NewCoapClient(":64868")
	client.Dial(":5683")


	req := NewRequest(TYPE_CONFIRMABLE, GET, 50782)
	req.SetRequestURI("/watch/this")
	req.Observe()

	resp, err := client.Send(req)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Got Synchronous Response:")
		log.Println(CoapCodeToString(resp.GetMessage().Code))
		log.Println(resp.GetMessage().Payload)
	}
}
