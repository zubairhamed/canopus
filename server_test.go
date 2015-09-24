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
	client := NewCoapServer(":0")
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
