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

	var notifyCount int = 0
	client.OnNotify(func (resource string, value interface{}, msg *Message) {
		// PrintMessage(msg)
		if notifyCount < 4 {
			notifyCount++
			log.Println("Got Change Notification for resource and value: ", notifyCount, resource, value)
		} else {
			log.Println("Cancelling Observation after 4 notifications")
			req := NewRequest(TYPE_CONFIRMABLE, GET, GenerateMessageId())
			req.SetRequestURI("watch/this")
			req.Observe(0)

			_, err := client.Send(req)
			if err != nil {
				log.Println(err)
			}
		}
	})

	client.Start()
}
