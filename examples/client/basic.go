package main

import (
	"github.com/zubairhamed/goap"
	"log"
)

/*	To test this example, also run examples/test_server.go */
func main() {
	// Client
	log.Println("Starting Client..")
	client := goap.NewClient()
	defer client.Close()

	client.Dial("udp", goap.COAP_DEFAULT_HOST, goap.COAP_DEFAULT_PORT)

	msg := goap.NewMessageOfType(goap.TYPE_CONFIRMABLE, 12345)
	msg.Code = goap.GET
	msg.Payload = []byte("Hello, goap")
	msg.AddOptions(goap.NewPathOptions("/example"))
	msg.Token = []byte(goap.GenerateToken(8))

	client.OnSuccess(func(msg *goap.Message) {
		log.Print("Got message back: " + string(msg.Payload))
	})

	client.OnReset(func(msg *goap.Message) {
		log.Println("Reset")
	})

	client.OnTimeout(func(msg *goap.Message) {
		log.Println("Timeout")
	})

	client.Send(msg)
}
