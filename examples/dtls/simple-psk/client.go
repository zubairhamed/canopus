package main

import (
	"fmt"

	"github.com/zubairhamed/canopus"
)

func main() {
	conn, err := canopus.DialDTLS("localhost:5684", "canopus", "secretPSK")
	if err != nil {
		panic(err.Error())
	}

	req := canopus.NewRequest(canopus.MessageConfirmable, canopus.Get)
	req.SetStringPayload("Hello, canopus")
	req.SetRequestURI("/hello")

	resp, err := conn.Send(req)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Got Response:" + resp.GetMessage().GetPayload().String())
}
