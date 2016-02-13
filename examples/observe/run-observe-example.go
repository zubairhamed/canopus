package main

import (
	"github.com/zubairhamed/canopus"
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
	server := canopus.NewLocalServer()
	server.Get("/watch/this", routeHandler)

	GenerateRandomChangeNotifications(server)

	server.OnMessage(func(msg *canopus.Message, inbound bool) {
		// PrintMessage(msg)
	})

	server.OnObserve(func(resource string, msg *canopus.Message) {
		log.Println("Observe Requested for " + resource)
	})

	server.Start()
}

func GenerateRandomChangeNotifications(server canopus.CoapServer) {
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

func routeHandler(req canopus.CoapRequest) canopus.CoapResponse {
	msg := canopus.NewMessageOfType(canopus.MessageAcknowledgment, req.GetMessage().MessageID)
	msg.SetStringPayload("Acknowledged")
	res := canopus.NewResponse(msg, nil)

	return res
}

func runClient() {
	client := canopus.NewCoapServer("0")

	client.OnStart(func(server canopus.CoapServer) {
		client.Dial("localhost:5683")
		req := canopus.NewRequest(canopus.MessageConfirmable, canopus.Get, canopus.GenerateMessageID())
		req.SetRequestURI("/watch/this")
		req.GetMessage().AddOption(canopus.OptionObserve, 0)

		_, err := client.Send(req)
		if err != nil {
			log.Println(err)
		}
	})

	var notifyCount int;
	client.OnNotify(func(resource string, value interface{}, msg *canopus.Message) {
		if notifyCount < 4 {
			notifyCount++
			log.Println("CLIENT: Got Change Notification for resource and value: ", notifyCount, resource, value)
		} else {
			log.Println("Cancelling Observation after 4 notifications")
			req := canopus.NewRequest(canopus.MessageConfirmable, canopus.Get, canopus.GenerateMessageID())
			req.SetRequestURI("watch/this")
			req.GetMessage().AddOption(canopus.OptionObserve, 0)

			_, err := client.Send(req)
			if err != nil {
				log.Println(err)
			}
		}
	})

	client.Start()
}
