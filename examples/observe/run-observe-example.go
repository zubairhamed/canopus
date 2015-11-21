package main

import (
	. "github.com/zubairhamed/canopus"
	"log"
	"math/rand"
	"strconv"
	"time"
)

func main() {
	go runServer()
	runClient()

	for {

	}
}

func runServer() {
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

func GenerateRandomChangeNotifications(server CoapServer) {
	ticker := time.NewTicker(3 * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				changeVal := strconv.Itoa(rand.Int())
				log.Println("SERVER: Notify Change..", changeVal)

				server.NotifyChange("/watch/this", changeVal, false)
			}
		}
	}()
}

func routeHandler(req CoapRequest) CoapResponse {
	msg := NewMessageOfType(TYPE_ACKNOWLEDGEMENT, req.GetMessage().MessageId)
	msg.SetStringPayload("Acknowledged")
	res := NewResponse(msg, nil)

	return res
}

func runClient() {
	client := NewCoapServer("0")

	client.OnStart(func(server CoapServer) {
		client.Dial("localhost:5683")
		req := NewRequest(TYPE_CONFIRMABLE, GET, GenerateMessageId())
		req.SetRequestURI("/watch/this")
		req.GetMessage().AddOption(OPTION_OBSERVE, 0)

		_, err := client.Send(req)
		if err != nil {
			log.Println(err)
		}
	})

	var notifyCount int = 0
	client.OnNotify(func(resource string, value interface{}, msg *Message) {
		if notifyCount < 4 {
			notifyCount++
			log.Println("CLIENT: Got Change Notification for resource and value: ", notifyCount, resource, value)
		} else {
			log.Println("Cancelling Observation after 4 notifications")
			req := NewRequest(TYPE_CONFIRMABLE, GET, GenerateMessageId())
			req.SetRequestURI("watch/this")
			req.GetMessage().AddOption(OPTION_OBSERVE, 0)

			_, err := client.Send(req)
			if err != nil {
				log.Println(err)
			}
		}
	})

	client.Start()
}
