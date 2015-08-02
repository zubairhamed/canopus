package main
import (
	. "github.com/zubairhamed/canopus"
	"log"
)

func main() {
	client := NewCoapServer(":0")

	client.OnStart(func(server *CoapServer) {
		client.Dial("localhost:5683")

		req := NewRequest(TYPE_CONFIRMABLE, GET, GenerateMessageId())
		req.SetProxyUri("http://api.openweathermap.org/data/2.1/find/city?bbox=12,32,15,37,10&cluster=yes")

		PrintMessage(req.GetMessage())
		resp, err := client.Send(req)

		log.Print(err)
		log.Print(resp)
	})

//	client.OnMessage(func(msg *Message, inbound bool) {
//		if inbound {
//			log.Println(">>>>> INBOUND <<<<<")
//		} else {
//			log.Println(">>>>> OUTBOUND <<<<<")
//		}
//
//		PrintMessage(msg)
//	})
	client.Start()
}
