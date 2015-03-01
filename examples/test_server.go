package main

import (
    "github.com/zubairhamed/goap"
)

func main() {
	server := goap.NewServer("udp", goap.COAP_DEFAULT_HOST)

	server.NewRoute("example", handleGet, goap.GET)
	server.NewRoute("example", handleDelete, goap.DELETE)
	server.NewRoute("example", handlePost, goap.POST)
	server.NewRoute("example", handlePut, goap.PUT)

    server.Start()
}

func handleGet(msg *goap.Message) *goap.Message {
    return createStandardResponse(msg)
}

func handlePost(msg *goap.Message) *goap.Message {
	return createStandardResponse(msg)
}

func handlePut(msg *goap.Message) *goap.Message {
	return createStandardResponse(msg)
}

func handleDelete(msg *goap.Message) *goap.Message {
	return createStandardResponse(msg)
}

func createStandardResponse(msg *goap.Message) *goap.Message {
	ack := goap.NewAcknowledgementMessage(msg.MessageId)
	ack.Payload = []byte("Hello GoAP")
	ack.Token = msg.Token
	ack.Code = goap.COAPCODE_205_CONTENT

	ack.AddOption(goap.NewOption(goap.OPTION_CONTENT_FORMAT, goap.MEDIATYPE_TEXT_PLAIN))

	return ack
}
