package goap

import (
	"sync"
	"testing"
)

func TestServer(t *testing.T) {
	// Done channel
	var wg sync.WaitGroup

	server := NewLocalServer()

	wg.Add(1)
	server.NewRoute("testGET", GET, func(msg *Message) *Message {
		defer wg.Done()
		if msg.MessageType != TYPE_CONFIRMABLE {
			t.Error("Type of message received should be TYPE_CONFIRMABLE")
		}

		// Payload
		if string(msg.Payload) != "TestGET" {
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

		if msg.GetPath() != "testGET" {
			t.Error("Path should be testGET")
		}

		if msg.MethodString() != "GET" {
			t.Error("MethodString should be GET")
		}

		if msg.MessageId != 12345 {
			t.Error("Message ID expected == 12345")
		}

		return msg
	})

	server.NewRoute("discoveryService1", GET, func(msg *Message) *Message { return msg })
	server.NewRoute("discoveryService2", GET, func(msg *Message) *Message { return msg })
	server.NewRoute("discoveryService3", GET, func(msg *Message) *Message { return msg })

	go server.Start()

	client := newTestClient()

	msg := NewMessageOfType(TYPE_CONFIRMABLE, 12345)
	msg.Code = GET
	msg.Payload = []byte("TestGET")
	msg.AddOptions(NewPathOptions("/testGET"))
	msg.Token = []byte(GenerateToken(8))

	go client.Send(msg)

	wg.Wait()
	client.Close()
}

func newTestClient() *Client {
	client := NewClient()
	client.Dial("udp", COAP_DEFAULT_HOST, COAP_DEFAULT_PORT)

	return client
}
