package main

import "github.com/zubairhamed/canopus"

func main() {
	server := canopus.NewServer()

	server.Get("/proxycall", func(req canopus.Request) canopus.Response {
		canopus.PrintMessage(req.GetMessage())
		msg := canopus.ContentMessage(req.GetMessage().GetMessageId(), canopus.MessageAcknowledgment)
		msg.SetStringPayload("Data from :5685 -- " + req.GetMessage().GetPayload().String())
		res := canopus.NewResponse(msg, nil)

		return res
	})
	server.ListenAndServe(":5685")
	<-make(chan struct{})
}
