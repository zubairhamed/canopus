package main
import (
	. "github.com/zubairhamed/canopus"
	"github.com/zubairhamed/go-commons/network"
	"time"
	"fmt"
)

func main() {
	server := NewLocalServer()
	server.NewRoute("/watch/this", GET, routeHandler)

	go GenerateRandomChangeNotifications()

	server.Start()
}

func GenerateRandomChangeNotifications() {
	for {
		time.Sleep(2 * time.Second)
		fmt.Println("bg!!")
	}
}


func routeHandler(r network.Request) network.Response {
	req := r.(*CoapRequest)
	msg := NewMessageOfType(TYPE_ACKNOWLEDGEMENT, req.GetMessage().MessageId)
	msg.SetStringPayload("Acknowledged")
	res := NewResponse(msg, nil)

	// server := req.GetServer()

	// server.Notify("/watch/this", )

	return res
}