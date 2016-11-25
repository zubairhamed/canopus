package main

import "github.com/zubairhamed/canopus"

func main() {
	server := canopus.NewServer()
	server.ProxyOverCoap(true)

	server.Get("/proxycall", func(req canopus.Request) canopus.Response {
		canopus.PrintMessage(req.GetMessage())
		msg := canopus.ContentMessage(req.GetMessage().GetMessageId(), canopus.MessageAcknowledgment)
		msg.SetStringPayload("Acknowledged: " + req.GetMessage().GetPayload().String())
		res := canopus.NewResponse(msg, nil)

		return res
	})
	server.ListenAndServe(":5683")
	<-make(chan struct{})
}
