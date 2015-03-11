package main

import (
	. "github.com/zubairhamed/goap"
	"log"
)

/*	To test this example, also run examples/test_server.go */
func main() {
	// Client
	log.Println("Starting Client..")
	client := NewClient()
	defer client.Close()

	client.Dial("udp", COAP_DEFAULT_HOST, COAP_DEFAULT_PORT)

	msg := NewMessageOfType(TYPE_CONFIRMABLE, 12345)
	msg.Code = GET
	msg.Payload = []byte("Hello, goap")
	msg.AddOptions(NewPathOptions("/example"))
	msg.Token = []byte(GenerateToken(8))

	client.OnSuccess(func(msg *Message) {
		log.Print("Got message back: " + string(msg.Payload))
	})

	client.OnReset(func(msg *Message) {
		log.Println("Reset")
	})

	client.OnTimeout(func(msg *Message) {
		log.Println("Timeout")
	})

	client.Send(msg)
}
