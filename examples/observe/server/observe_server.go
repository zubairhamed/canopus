package main
import (
	. "github.com/zubairhamed/canopus"
	"github.com/zubairhamed/go-commons/network"
)

func main() {
	server := NewLocalServer()
	server.NewRoute("/watch/this", GET, routeHandler)

	server.Start()
}


func routeHandler(r network.Request) network.Response {
	req := r.(*CoapRequest)
	msg := NewMessageOfType(TYPE_ACKNOWLEDGEMENT, req.GetMessage().MessageId)
	msg.SetStringPayload("Acknowledged")
	res := NewResponse(msg, nil)

	return res
}