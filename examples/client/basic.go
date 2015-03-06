package main

import (
	"github.com/zubairhamed/goap"
	"log"
)

func main() {
	// Client
	log.Println("Starting Client..")
	client := goap.NewClient()
	client.Dial("udp", goap.COAP_DEFAULT_HOST)

	msg := goap.NewMessageOfType(goap.TYPE_CONFIRMABLE, 12345)
	msg.Code = goap.GET
	msg.Payload = []byte("Hello, goap")
	msg.AddOptions(goap.NewPathOptions("/example"))

	client.OnSuccess(func(msg *goap.Message) {
		log.Println("Success")
	})

	client.OnReset(func(msg *goap.Message) {
		log.Println("Reset")
	})

	client.OnTimeout(func(msg *goap.Message){
		log.Println("Timeout")
	})

	client.Send(msg)
}

