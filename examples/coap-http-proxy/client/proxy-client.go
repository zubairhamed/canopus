package main

import (
	. "github.com/zubairhamed/canopus"
)

func main() {
	client := NewCoapServer(":0")

	client.OnStart(func(server *CoapServer) {
		client.Dial("localhost:5683")

		req := NewRequest(TYPE_CONFIRMABLE, GET, GenerateMessageId())
		req.SetProxyUri("http://api.openweathermap.org/data/2.5/weather?q=London,uk&appid=2de143494c0b295cca9337e1e96b00e0")

		PrintMessage(req.GetMessage())
		resp, err := client.Send(req)
		if err != nil {
			println("err", err)
		}
		PrintMessage(resp.GetMessage())
	})
	client.Start()
}
