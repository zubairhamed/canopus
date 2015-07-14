package main
import (
	. "github.com/zubairhamed/canopus"
	"github.com/zubairhamed/go-commons/network"
	"time"
	"log"
	"math/rand"
)

func main() {
	server := NewLocalServer()
	server.NewRoute("watch/this", GET, routeHandler)

	GenerateRandomChangeNotifications(server)

	server.OnMessage(func (msg *Message, inbound bool){
		// PrintMessage(msg)
	})

	server.OnObserve(func(resource string, msg *Message){
		log.Println("Observe Requested for " + resource)
	})

	server.Start()
}

func GenerateRandomChangeNotifications(server *CoapServer) {
	ticker := time.NewTicker(3 * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				log.Println("Notify Change..")
				log.Println(rand.Float32())

				server.NotifyChange("watch/this", "Some new value", false)
			}
		}
	}()
}

func routeHandler(r network.Request) network.Response {
	req := r.(*CoapRequest)
	msg := NewMessageOfType(TYPE_ACKNOWLEDGEMENT, req.GetMessage().MessageId)
	msg.SetStringPayload("Acknowledged")
	res := NewResponse(msg, nil)

	return res
}