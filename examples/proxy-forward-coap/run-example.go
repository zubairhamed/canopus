package main

import (
	"github.com/zubairhamed/canopus"
)

func main() {
	go runProxyServer()
	go runServer()
	go runClient()

	<-make(chan struct{})
}

func runProxyServer() {
	server := canopus.NewLocalServer("TestServer")
	server.ProxyCoap(true)

	server.Get("/proxycall", func(req canopus.CoapRequest) canopus.CoapResponse {
		canopus.PrintMessage(req.GetMessage())
		msg := canopus.ContentMessage(req.GetMessage().MessageID, canopus.MessageAcknowledgment)
		msg.SetStringPayload("Acknowledged: " + req.GetMessage().Payload.String())
		res := canopus.NewResponse(msg, nil)

		return res
	})
	server.Start()
}

func runServer() {
	server := canopus.NewCoapServer("TestServer", ":5684")

	server.Get("/proxycall", func(req canopus.CoapRequest) canopus.CoapResponse {
		canopus.PrintMessage(req.GetMessage())
		msg := canopus.ContentMessage(req.GetMessage().MessageID, canopus.MessageAcknowledgment)
		msg.SetStringPayload("Data from :5684 -- " + req.GetMessage().Payload.String())
		res := canopus.NewResponse(msg, nil)

		return res
	})

	server.Start()
}

func runClient() {
	client := canopus.NewCoapServer("TestServer", ":0")

	client.OnStart(func(server canopus.CoapServer) {
		client.Dial("localhost:5683")

		req := canopus.NewRequest(canopus.MessageConfirmable, canopus.Get, canopus.GenerateMessageID())
		req.SetProxyURI("coap://localhost:5684/proxycall")

		canopus.PrintMessage(req.GetMessage())
		resp, err := client.Send(req)
		if err != nil {
			println("err", err)
		}
		canopus.PrintMessage(resp.GetMessage())
	})
	client.Start()
}
