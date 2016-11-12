package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/zubairhamed/canopus"
)

func main() {
	go runServer()
	go runClient()

	<-make(chan struct{})
}

func runServer() {
	server := canopus.NewServer()
	server.Get("/watch/this", func(req canopus.Request) canopus.Response {
		msg := canopus.NewMessageOfType(canopus.MessageAcknowledgment, req.GetMessage().GetMessageId(), canopus.NewPlainTextPayload("Acknowledged"))
		res := canopus.NewResponse(msg, nil)

		return res
	})

	ticker := time.NewTicker(3 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				changeVal := strconv.Itoa(rand.Int())
				fmt.Println("[SERVER << ] Change of value -->", changeVal)

				server.NotifyChange("/watch/this", changeVal, false)
			}
		}
	}()

	server.OnMessage(func(msg canopus.Message, inbound bool) {
		canopus.PrintMessage(msg)
	})

	server.OnObserve(func(resource string, msg canopus.Message) {
		fmt.Println("[SERVER << ] Observe Requested for " + resource)
	})

	server.ListenAndServe(":5683")
}

func runClient() {
	client := canopus.NewClient()
	conn, err := client.Dial("localhost:5683")

	tok, err := conn.ObserveResource("/watch/this")
	if err != nil {
		panic(err.Error())
	}

	obsChannel := make(chan canopus.ObserveMessage)
	done := make(chan bool)
	go conn.Observe(obsChannel)

	notifyCount := 0
	go func() {
		for {
			select {
			case obsMsg, open := <-obsChannel:
				if open {
					if notifyCount == 5 {
						fmt.Println("[CLIENT >> ] Canceling observe after 5 notifications..")
						go conn.CancelObserveResource("watch/this", tok)
						go conn.StopObserve(obsChannel)
						done <- true
						return
					} else {
						notifyCount++
						// msg := obsMsg.Msg\
						resource := obsMsg.GetResource()
						val := obsMsg.GetValue()

						fmt.Println("[CLIENT >> ] Got Change Notification for resource and value: ", notifyCount, resource, val)
					}
				} else {
					done <- true
					return
				}
			}
		}
	}()
	<-done
	fmt.Println("Done")
}
