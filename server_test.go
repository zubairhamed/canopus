package canopus

import (
	"sync"
	"testing"
	"github.com/zubairhamed/go-commons/network"
	"fmt"
)

func TestServer(t *testing.T) {
	// Done channel
	var wg sync.WaitGroup

	server := NewLocalServer()

	wg.Add(1)

	server.NewRoute("testGET", GET, func(req network.Request) network.Response {
		msg := req.(*CoapRequest).GetMessage()

		defer wg.Done()

		if msg.MessageType != TYPE_CONFIRMABLE {
			t.Error("Type of message received should be TYPE_CONFIRMABLE")
		}

		// Payload
		if string(msg.Payload.String()) != "TestGET" {
			t.Error("Payload should be TestGET")
		}

		// Method
		if msg.Code != GET {
			t.Error("Message Code/Method should be GET")
		}

		if msg.Token == nil {
			t.Error("Token should not be a nil value")
		}

		if msg.GetTokenLength() != 8 {
			t.Error("Token length should be 8")
		}

		if msg.GetUriPath() != "testGET" {
			t.Error("Path should be testGET")
		}

		if msg.MethodString() != "GET" {
			t.Error("MethodString should be GET")
		}

		if msg.MessageId != 12345 {
			t.Error("Message ID expected == 12345")
		}

		return NewResponse(msg, nil)
	})

	server.NewRoute("discoveryService1", GET, func(req network.Request) network.Response { return NewResponse(req.(*CoapRequest).GetMessage(), nil) })
	server.NewRoute("discoveryService2", GET, func(req network.Request) network.Response { return NewResponse(req.(*CoapRequest).GetMessage(), nil) })
	server.NewRoute("discoveryService3", GET, func(req network.Request) network.Response { return NewResponse(req.(*CoapRequest).GetMessage(), nil) })

	go server.Start()

	client := newTestClient()
	defer client.Close()

	msg := NewMessageOfType(TYPE_CONFIRMABLE, 12345)
	msg.Code = GET
	msg.Payload = network.NewPlainTextPayload("TestGET")
	msg.AddOptions(NewPathOptions("/testGET"))
	msg.Token = []byte(GenerateToken(8))

	go client.Send(NewRequestFromMessage(msg))

	wg.Wait()
}

func newTestClient() *CoapClient {
	client := NewCoapClient(":12345")

	client.Dial(fmt.Sprintf("%s:%d", COAP_DEFAULT_HOST, COAP_DEFAULT_PORT))

	return client
}
