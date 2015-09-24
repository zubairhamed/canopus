package canopus

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestServerInstantiate(t *testing.T) {
	var s *CoapServer
	s = NewCoapServer("localhost:1000")
	assert.Equal(t, 1000, s.localAddr.Port)
	assert.Equal(t, "udp", s.localAddr.Network())

	s = NewLocalServer()
	assert.NotNil(t, s)
	assert.Equal(t, 5683, s.localAddr.Port)
	assert.Equal(t, "udp", s.localAddr.Network())
}

func TestDiscoveryService(t *testing.T) {
	server := NewLocalServer()
	assert.NotNil(t, server)
	assert.Equal(t, 5683, server.localAddr.Port)
	assert.Equal(t, "udp", server.localAddr.Network())

	go server.Start()
	client := NewCoapClient()
	client.OnStart(func(server *CoapServer) {
		tok := "abc123"
		client.Dial("localhost:5683")

		req := NewRequest(TYPE_CONFIRMABLE, GET, GenerateMessageId())
		req.SetToken(tok)
		req.SetRequestURI(".well-known/core")
		resp, err := client.Send(req)
		assert.Nil(t, err)

		assert.Equal(t, tok, resp.GetMessage().GetTokenString())
		client.Stop()
	})
	client.Start()
}

func TestClientServerRequestResponse(t *testing.T) {
	server := NewLocalServer()

	server.Get("/ep", func (req CoapRequest) CoapResponse {
		msg := ContentMessage(req.GetMessage().MessageId, TYPE_ACKNOWLEDGEMENT)
		msg.SetStringPayload("ACK GET")
		res := NewResponse(msg, nil)

		return res
	})
	go server.Start()

	client := NewCoapClient()

	client.OnStart(func(server *CoapServer) {
		client.Dial("localhost:5683")
		token := "tok1234"

		var req CoapRequest
		var resp CoapResponse
		var err error

		// 404 Test
		req = NewConfirmableGetRequest()
		req.SetToken(token)
		req.SetRequestURI("ep-404")
		resp, err = client.Send(req)

		assert.Equal(t, COAPCODE_404_NOT_FOUND, resp.GetMessage().Code)

		// Test for Token
		req = NewConfirmableGetRequest()
		req.SetToken(token)
		req.SetRequestURI("ep")
		resp, err = client.Send(req)

		assert.Nil(t, err)
		assert.Equal(t, "ACK GET", resp.GetMessage().Payload.String())
		assert.Equal(t, token, resp.GetMessage().GetTokenString())

		// Test default token set
		req = NewConfirmableGetRequest()
		req.SetRequestURI("ep")
		resp, err = client.Send(req)

		assert.Nil(t, err)
		assert.Equal(t, "ACK GET", resp.GetMessage().Payload.String())
		assert.NotEmpty(t, resp.GetMessage().GetTokenString())

		client.Stop()
	})
	client.Start()
}
