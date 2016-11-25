package main

import (
	"fmt"

	"github.com/zubairhamed/canopus"
)

func main() {
	conn, err := canopus.DialDTLS("localhost:5684", "canopus", "secretPSK")
	if err != nil {
		panic(err.Error())
	}

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
