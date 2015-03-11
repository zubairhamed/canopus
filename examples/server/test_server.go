package main

import (
	. "github.com/zubairhamed/goap"
	"log"
)

func main() {
	log.Println("Starting up server..")
	server := NewLocalServer()

	server.NewRoute("example", GET, func(msg *Message) *Message {
		return createStandardResponse(msg)
	})

	server.NewRoute("example", DELETE, func(msg *Message) *Message {
		return createStandardResponse(msg)
	})

	server.NewRoute("example", POST, func(msg *Message) *Message {
		return createStandardResponse(msg)
	})

	server.NewRoute("example", PUT, func(msg *Message) *Message {
		return createStandardResponse(msg)
	})
	log.Println("Server Started")
	server.Start()
}

func createStandardResponse(msg *Message) *Message {
	log.Println("Got Request ")
	ack := NewMessageOfType(TYPE_ACKNOWLEDGEMENT, msg.MessageId)
	// ack := NewMessageOfType(TYPE_ACKNOWLEDGEMENT, binary.BigEndian.Uint16([]byte{10, 20}))
	ack.Code = COAPCODE_205_CONTENT
	ack.Token = msg.Token
	ack.Payload = []byte("Hello GoAP")

	ack.AddOption(OPTION_CONTENT_FORMAT, MEDIATYPE_APPLICATION_XML)

	return ack
}
