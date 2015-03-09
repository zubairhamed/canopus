package main

import (
	"github.com/zubairhamed/goap"
	// "encoding/binary"
	"log"
)

func main() {
	log.Println("Starting up server..")
	server := goap.NewLocalServer()

	server.NewRoute("example", goap.GET, func(msg *goap.Message) *goap.Message {
		return createStandardResponse(msg)
	})

	server.NewRoute("example", goap.DELETE, func(msg *goap.Message) *goap.Message {
		return createStandardResponse(msg)
	})

	server.NewRoute("example", goap.POST, func(msg *goap.Message) *goap.Message {
		return createStandardResponse(msg)
	})

	server.NewRoute("example", goap.PUT, func(msg *goap.Message) *goap.Message {
		return createStandardResponse(msg)
	})
	log.Println("Server Started")
	server.Start()
}

func createStandardResponse(msg *goap.Message) *goap.Message {
	log.Println("Got Request ")
	ack := goap.NewMessageOfType(goap.TYPE_ACKNOWLEDGEMENT, msg.MessageId)
	// ack := goap.NewMessageOfType(goap.TYPE_ACKNOWLEDGEMENT, binary.BigEndian.Uint16([]byte{10, 20}))
	ack.Code = goap.COAPCODE_205_CONTENT
	ack.Token = msg.Token
	ack.Payload = []byte("Hello GoAP")

	ack.AddOption(goap.OPTION_CONTENT_FORMAT, goap.MEDIATYPE_APPLICATION_XML)

	return ack
}
