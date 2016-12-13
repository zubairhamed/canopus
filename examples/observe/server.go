package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/zubairhamed/canopus"
)

func main() {
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
	<-make(chan struct{})
}
