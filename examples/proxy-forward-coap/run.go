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
	server := canopus.NewServer()
	server.ProxyOverCoap(true)

	server.Get("/proxycall", func(req canopus.Request) canopus.Response {
		canopus.PrintMessage(req.GetMessage())
		msg := canopus.ContentMessage(req.GetMessage().GetMessageId(), canopus.MessageAcknowledgment)
		msg.SetStringPayload("Acknowledged: " + req.GetMessage().GetPayload().String())
		res := canopus.NewResponse(msg, nil)

		return res
	})
	server.ListenAndServe(":5683", nil)
}

func runServer() {
	server := canopus.NewServer()

	server.Get("/proxycall", func(req canopus.Request) canopus.Response {
		canopus.PrintMessage(req.GetMessage())
		msg := canopus.ContentMessage(req.GetMessage().GetMessageId(), canopus.MessageAcknowledgment)
		msg.SetStringPayload("Data from :5685 -- " + req.GetMessage().GetPayload().String())
		res := canopus.NewResponse(msg, nil)

		return res
	})

	server.ListenAndServe(":5685", nil)
}

func runClient() {
	client := canopus.NewClient()
	conn, err := client.Dial("localhost:5683")
	if err != nil {
		panic(err.Error())
	}

	req := canopus.NewRequest(canopus.MessageConfirmable, canopus.Get, canopus.GenerateMessageID())
	req.SetProxyURI("coap://localhost:5685/proxycall")

	canopus.PrintMessage(req.GetMessage())
	resp, err := conn.Send(req)
	if err != nil {
		println("err", err)
	}
	canopus.PrintMessage(resp.GetMessage())
}
