package main

import (
	"fmt"
	"log"
	"time"

	"github.com/zubairhamed/canopus"
)

func main() {
	fmt.Println("Connecting to CoAP Server")
	conn, err := canopus.Dial("localhost:5683")
	if err != nil {
		panic(err.Error())
	}
	//
	//req := canopus.NewRequest(canopus.MessageConfirmable, canopus.Get, canopus.GenerateMessageID()).(*canopus.CoapRequest)
	//req.SetStringPayload("Hello, canopus")
	//req.SetRequestURI("/hello")
	//
	//fmt.Println("Sending request..")
	//resp, err := conn.Send(req)
	//if err != nil {
	//	panic(err.Error())
	//}
	//
	//fmt.Println("Got Response:" + resp.GetMessage().GetPayload().String())

	t := time.NewTicker(time.Second)
	for {
		log.Println("A1")
		req := canopus.NewRequest(canopus.MessageConfirmable, canopus.Get, canopus.GenerateMessageID())
		log.Println("A2")
		req.SetStringPayload("Hello, canopus")
		log.Println("A3")
		req.SetRequestURI("/hello")
		log.Println("A4")
		resp, err := conn.Send(req)
		log.Println("A5")
		log.Println(resp, err)
		log.Println("A6")
		if err != nil {
			log.Println("A7")
			panic(err.Error())
		}
		log.Println("A8")
		log.Println("Got Response:" + resp.GetMessage().GetPayload().String())
		<-t.C
		log.Println("A9")
	}
}
