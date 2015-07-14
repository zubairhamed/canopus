package canopus

//import (
//	"sync"
//	"testing"
//	"log"
//)
//
//func TestServer(t *testing.T) {
//	// Done channel
//	var wg sync.WaitGroup
//
//	server := NewLocalServer()
//
//	wg.Add(1)
//
//	server.NewRoute("testGET", GET, func(req network.Request) network.Response {
//		msg := req.(*CoapRequest).GetMessage()
//
//		defer wg.Done()
//
//		if msg.MessageType != TYPE_CONFIRMABLE {
//			t.Error("Type of message received should be TYPE_CONFIRMABLE")
//		}
//
//		// Payload
//		if string(msg.Payload.String()) != "TestGET" {
//			t.Error("Payload should be TestGET")
//		}
//
//		// Method
//		if msg.Code != GET {
//			t.Error("Message Code/Method should be GET")
//		}
//
//		if msg.Token == nil {
//			t.Error("Token should not be a nil value")
//		}
//
//		if msg.GetTokenLength() != 8 {
//			t.Error("Token length should be 8")
//		}
//
//		if msg.GetUriPath() != "testGET" {
//			t.Error("Path should be testGET")
//		}
//
//		if msg.MethodString() != "GET" {
//			t.Error("MethodString should be GET")
//		}
//
//		if msg.MessageId != 12345 {
//			t.Error("Message ID expected == 12345")
//		}
//
//		return NewResponse(msg, nil)
//	})
//
//	server.NewRoute("discoveryService1", GET, func(req network.Request) network.Response { return NewResponse(req.(*CoapRequest).GetMessage(), nil) })
//	server.NewRoute("discoveryService2", GET, func(req network.Request) network.Response { return NewResponse(req.(*CoapRequest).GetMessage(), nil) })
//	server.NewRoute("discoveryService3", GET, func(req network.Request) network.Response { return NewResponse(req.(*CoapRequest).GetMessage(), nil) })
//
//	server.OnStart(func(s *CoapServer){
//		client := NewCoapServer("localhost:0")
//
//		msg := NewMessageOfType(TYPE_CONFIRMABLE, 12345)
//		msg.Code = GET
//		msg.Payload = network.NewPlainTextPayload("TestGET")
//		msg.AddOptions(NewPathOptions("/testGET"))
//		msg.Token = []byte(GenerateToken(8))
//
//		client.Dial("localhost:5683")
//		client.OnStart(func(s *CoapServer){
//			log.Println(".....")
//			resp, err := client.Send(NewRequestFromMessage(msg))
//			log.Println(".....2")
//			log.Println(err)
//			PrintMessage(resp.GetMessage())
//		})
//		client.Start()
//	})
//	server.Start()
//
//	// wg.Wait()
//}
