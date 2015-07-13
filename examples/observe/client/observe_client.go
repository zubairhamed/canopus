package main

import (
	. "github.com/zubairhamed/canopus"
	"log"
 	"net"
)

func main() {
	client := NewCoapServer(":0")

	client.On(EVT_START, func() {
		client.Dial("localhost:5683")
		req := NewRequest(TYPE_CONFIRMABLE, GET, GenerateMessageId())
		req.SetRequestURI("watch/this")
		req.Observe()

		addr, err := net.ResolveUDPAddr("udp", "localhost:5683")
		if err != nil {
			log.Println(err)
		}

		resp, err := client.ServerSendTo(req, addr)
		if err != nil {
			log.Println(err)
		} else {
			log.Println("Got Synchronous Response:")
			log.Println(CoapCodeToString(resp.GetMessage().Code))
			log.Println(resp.GetMessage().Payload)
		}
	})

	client.On(EVT_NOTIFY, func(){

	})

	client.Start()
}
