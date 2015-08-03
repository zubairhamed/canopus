package main

import (
	. "github.com/zubairhamed/canopus"
)

func main() {
	client := NewCoapServer(":0")

	client.OnStart(func(server *CoapServer) {
		client.Dial("localhost:5683")

		req := NewRequest(TYPE_CONFIRMABLE, GET, GenerateMessageId())
		req.SetProxyUri("http://api.openweathermap.org/data/2.1/find/city?bbox=12,32,15,37,10&cluster=yes")

		PrintMessage(req.GetMessage())
		resp, _ := client.Send(req)

		PrintMessage(resp.GetMessage())
	})
	client.Start()
}
