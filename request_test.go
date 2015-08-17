package canopus

import (
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func TestRequest(t *testing.T) {
	var req CoapRequest
	assert.NotNil(t, NewRequest(TYPE_CONFIRMABLE, GET, 12345))

	msg := NewMessage(TYPE_CONFIRMABLE, GET, 12345)

	req = NewRequestFromMessage(msg)
	assert.NotNil(t, req)

	assert.Equal(t, GET, req.GetMessage().Code)
	assert.NotNil(t, NewClientRequestFromMessage(msg, make(map[string]string), &net.UDPConn{}, &net.UDPAddr{}))

	req = NewRequest(TYPE_CONFIRMABLE, GET, 12345)
	assert.Equal(t, uint8(0), req.GetMessage().MessageType)

	req.SetConfirmable(false)
	assert.Equal(t, uint8(1), req.GetMessage().MessageType)
}
