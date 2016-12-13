package main

import (
	"fmt"

	"github.com/zubairhamed/canopus"
)

func main() {
	fmt.Println("Connecting to CoAP Server")
	conn, err := canopus.Dial("localhost:5683")
	if err != nil {
		panic(err.Error())
	}

	req := canopus.NewRequest(canopus.MessageConfirmable, canopus.Get, canopus.GenerateMessageID()).(*canopus.CoapRequest)
	req.SetStringPayload("Hello, canopus")
	req.SetRequestURI("/hello")

	fmt.Println("Sending request..")
	resp, err := conn.Send(req)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Got Response:" + resp.GetMessage().GetPayload().String())
}
