package main
import (
	. "github.com/zubairhamed/canopus"
	"github.com/zubairhamed/go-commons/network"
	"log"
	"time"
)

func main() {
	server := NewLocalServer()
	server.NewRoute("/watch/this", GET, routeHandler)

	go server.Start()

	time.Sleep(3 * time.Second)

	client := NewCoapClient(":64868")
	client.Dial(":5683")

	req := NewRequest(TYPE_CONFIRMABLE, GET, 50782)
	req.SetRequestURI("/.well-known/core")

	// resp, err := client.Send(req)
	client.SendAsync(req, func(resp *CoapResponse, err error) {
		if err != nil {
			log.Println(err)
		} else {
			log.Println("Got Asynchronous Response:")
			log.Println(CoapCodeToString(resp.GetMessage().Code))
		}
	})


	// Observe Resource
//	req := NewRequest(TYPE_CONFIRMABLE, GET, 50783)
//	req.SetRequestURI("/watch/this")
//	req.Observe()

	for {}
}


func routeHandler(r network.Request) network.Response {
	req := r.(*CoapRequest)
	msg := NewMessageOfType(TYPE_ACKNOWLEDGEMENT, req.GetMessage().MessageId)
	msg.SetStringPayload("Acknowledged")
	res := NewResponse(msg, nil)

	return res
}