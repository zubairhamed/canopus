package main

import (
	. "github.com/zubairhamed/canopus"
	"log"
)

func main() {
	client := NewCoapServer(":0")

	client.OnStart(func (server *CoapServer){
		client.Dial("localhost:5683")
		req := NewRequest(TYPE_CONFIRMABLE, GET, GenerateMessageId())
		req.SetRequestURI("watch/this")
		req.Observe(0)

		_, err := client.Send(req)
		if err != nil {
			log.Println(err)
		}
	})

	client.OnNotify(func (resource string, value interface{}, msg *Message) {
		PrintMessage(msg)
		log.Println("Got Change Notification for resource and value: ", resource, value)
	})

	client.Start()
}
