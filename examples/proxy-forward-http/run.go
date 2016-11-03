package main

import (
	"github.com/zubairhamed/canopus"
)

func main() {
	go runServer()
	go runClient()

	<-make(chan struct{})
}

func runClient() {
	client := canopus.NewClient()
	conn, err := client.Dial("localhost:5683")
	if err != nil {
		panic(err.Error())
	}

	req := canopus.NewRequest(canopus.MessageConfirmable, canopus.Get, canopus.GenerateMessageID())
	req.SetProxyURI("http://api.openweathermap.org/data/2.5/weather?q=London,uk&appid=2de143494c0b295cca9337e1e96b00e0")

	canopus.PrintMessage(req.GetMessage())
	resp, err := conn.Send(req)
	if err != nil {
		println("err", err)
	}
	canopus.PrintMessage(resp.GetMessage())
}

func runServer() {
	server := canopus.NewServer()
	server.ProxyHTTP(true)

	server.ListenAndServe(":5683", nil)
}
