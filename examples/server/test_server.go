package main

import (
    "github.com/zubairhamed/goap"
)

func main() {
	server := goap.NewServer("udp", goap.COAP_DEFAULT_HOST)

	server.NewRoute("example", goap.GET, func(msg *goap.Message) *goap.Message {
        return createStandardResponse(msg)
    })

	server.NewRoute("example", goap.DELETE, func(msg *goap.Message) *goap.Message {
        return createStandardResponse(msg)
    })

	server.NewRoute("example", goap.POST, func(msg *goap.Message) *goap.Message {
        return createStandardResponse(msg)
    })

	server.NewRoute("example", goap.PUT, func(msg *goap.Message) *goap.Message  {
        return createStandardResponse(msg)
    })

    server.Start()
}

func createStandardResponse(msg *goap.Message) *goap.Message {
    ack := goap.NewMessageOfType(goap.TYPE_ACKNOWLEDGEMENT, msg.MessageId)
	ack.Payload = []byte("Hello GoAP")
	ack.Token = msg.Token
	ack.Code = goap.COAPCODE_205_CONTENT

	ack.AddOption(goap.NewOption(goap.OPTION_CONTENT_FORMAT, goap.MEDIATYPE_TEXT_PLAIN))

	return ack
}
