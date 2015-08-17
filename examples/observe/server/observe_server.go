package main

import (
	. "github.com/zubairhamed/canopus"
	"log"
	"math/rand"
	"strconv"
	"time"
)

func main() {
	server := NewLocalServer()
	server.Get("/watch/this", routeHandler)

	GenerateRandomChangeNotifications(server)

	server.OnMessage(func(msg *Message, inbound bool) {
		// PrintMessage(msg)
	})

	server.OnObserve(func(resource string, msg *Message) {
		log.Println("Observe Requested for " + resource)
	})

	server.Start()
}

func GenerateRandomChangeNotifications(server *CoapServer) {
	ticker := time.NewTicker(3 * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				changeVal := strconv.Itoa(rand.Int())
				log.Println("Notify Change..", changeVal)

				server.NotifyChange("/watch/this", changeVal, false)
			}
		}
	}()
}

func routeHandler(req *Request) Response {
	msg := NewMessageOfType(TYPE_ACKNOWLEDGEMENT, req.GetMessage().MessageId)
	msg.SetStringPayload("Acknowledged")
	res := NewResponse(msg, nil)

	return res
}
