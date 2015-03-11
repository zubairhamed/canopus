package goap

import (
	"sync"
	"testing"
    . "github.com/zubairhamed/goap"
)

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
		} else {

        }

		// Payload
		if string(msg.Payload) != "TestGET" {
			t.Error("Payload should be TestGET")
			t.Fail()
		}

		// Method
		if msg.Code != GET {
			t.Error("Message Code/Method should be GET")
			t.Fail()
		}

		if msg.Token == nil {
			t.Error("Token should not be a nil value")
			t.Fail()
		}

        if msg.GetTokenLength() != 8 {
            t.Error("Token length should be 8")
            t.Fail()
        }

		if msg.GetPath() != "testGET" {
			t.Error("Path should be testGET")
			t.Fail()
		}

        if msg.MethodString() != "GET" {
            t.Error("MethodString should be GET")
            t.Fail()
        }

        if msg.MessageId != 12345 {
            t.Error("Message ID expected == 12345")
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

func newTestClient() *Client {
	client := NewClient()
	client.Dial("udp", COAP_DEFAULT_HOST, COAP_DEFAULT_PORT)

	return client
}
