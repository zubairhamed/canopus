package main

import (
	. "github.com/zubairhamed/goap"
	"log"
)

func main() {
	log.Println("Starting up lwm2m..")
	server := NewLocalServer()

	server.NewRoute("example", GET, func(req *CoapRequest) *CoapResponse {
		return createStandardResponse(req)
	})

	server.NewRoute("example", DELETE, func(req *CoapRequest) *CoapResponse {
		return createStandardResponse(req)
	})

	server.NewRoute("example", POST, func(req *CoapRequest) *CoapResponse {
		return createStandardResponse(req)
	})

	server.NewRoute("example", PUT, func(req *CoapRequest) *CoapResponse {
		return createStandardResponse(req)
	})
	log.Println("Server Started")
	server.Start()
}

func createStandardResponse(req *CoapRequest) *CoapResponse {
	log.Println("Got Request ")
	msg := req.GetMessage()
	ack := NewMessageOfType(TYPE_ACKNOWLEDGEMENT, msg.MessageId)
	// ack := NewMessageOfType(TYPE_ACKNOWLEDGEMENT, binary.BigEndian.Uint16([]byte{10, 20}))
	ack.Code = COAPCODE_205_CONTENT
	ack.Token = msg.Token
	ack.Payload = []byte("Hello GoAP")

	ack.AddOption(OPTION_CONTENT_FORMAT, MEDIATYPE_APPLICATION_XML)

	resp := NewResponse(ack, nil)

	return resp
}
