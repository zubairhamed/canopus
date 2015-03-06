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

	msg := goap.NewMessageOfType(goap.TYPE_CONFIRMABLE)
	msg.Code = goap.GET
	msg.MessageId = 12345
	msg.Payload = []byte{"hello, goap"}

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

func createStandardResponse(msg *goap.Message) *goap.Message {
	ack := goap.NewMessageOfType(goap.TYPE_ACKNOWLEDGEMENT, msg.MessageId)
	// ack := goap.NewMessageOfType(goap.TYPE_ACKNOWLEDGEMENT, binary.BigEndian.Uint16([]byte{10, 20}))
	ack.Code = goap.COAPCODE_205_CONTENT
	ack.Token = msg.Token
	ack.Payload = []byte("Hello GoAP")

	ack.AddOption(goap.OPTION_CONTENT_FORMAT, goap.MEDIATYPE_APPLICATION_XML)

	return ack
}
