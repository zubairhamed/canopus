package main

import "github.com/zubairhamed/canopus"

func main() {
	conn, err := canopus.Dial("localhost:5683")
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
