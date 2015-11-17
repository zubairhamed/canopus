package main

import (
	. "github.com/zubairhamed/canopus"
)

func main() {
	go runProxyServer()
	go runServer()
	runClient()

	for {

	}
}

func runProxyServer() {
	server := NewLocalServer()
	server.SetProxy(PROXY_COAP_COAP, true)

	server.Start()
}

func runServer() {
	server := NewCoapServer(":5684")

	server.Get("/proxycall", func(req CoapRequest) CoapResponse {
		PrintMessage(req.GetMessage())
		msg := ContentMessage(req.GetMessage().MessageId, TYPE_ACKNOWLEDGEMENT)
		msg.SetStringPayload("Acknowledged: " + req.GetMessage().Payload.String())
		res := NewResponse(msg, nil)

		return res
	})

	server.Start()
}

func runClient() {
	client := NewCoapServer(":0")

	client.OnStart(func(server *CoapServer) {
		client.Dial("localhost:5683")

		req := NewRequest(TYPE_CONFIRMABLE, GET, GenerateMessageId())
		req.SetProxyUri("coap://localhost:5684/proxycall")

		PrintMessage(req.GetMessage())
		resp, err := client.Send(req)
		if err != nil {
			println("err", err)
		}
		PrintMessage(resp.GetMessage())
	})
	client.Start()
}