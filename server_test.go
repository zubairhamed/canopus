package goap

import (
	"testing"
	"sync"
)


/*
	- Spin up test server as goroutine
	- Setup Channel
	- Run Tests
	-

 */

func TestServer(t *testing.T) {

	// Done channel
	var wg sync.WaitGroup

	server := NewLocalServer()

	wg.Add(1)
	server.NewRoute("testGET", GET, func(msg *Message) *Message {
		defer wg.Done()
		if msg.MessageType != TYPE_CONFIRMABLE {
			t.Log("Type of message received should be TYPE_CONFIRMABLE")
			t.Fail()
		}

		// Payload
		if string(msg.Payload) != "TestGET" {
			t.Log("Payload should be TestGET")
			t.Fail()
		}

		// Method
		if msg.Code != GET {
			t.Log("Message Code/Method should be GET")
			t.Fail()
		}

		if msg.Token == nil {
			t.Log("Token should not be a nil value")
			t.Fail()
		}

		if msg.GetPath() != "testGET" {
			t.Log("Path should be testGET")
			t.Fail()
		}

		return msg
	})

	go server.Start()

	client := newTestClient()
	defer client.Close()

	msg := NewMessageOfType(TYPE_CONFIRMABLE, 12345)
	msg.Code = GET
	msg.Payload = []byte("TestGET")
	msg.AddOptions(NewPathOptions("/testGET"))
	msg.Token = []byte(GenerateToken(8))

	go client.Send(msg)

	wg.Wait()
}

func newTestClient() (*Client) {
	client := NewClient()
	client.Dial("udp", COAP_DEFAULT_HOST, COAP_DEFAULT_PORT)

	return client
}
