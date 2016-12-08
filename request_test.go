package canopus

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequest(t *testing.T) {
	var req Request
	assert.NotNil(t, NewRequest(MessageConfirmable, Get, 12345))

	msg := NewMessage(MessageConfirmable, Get, 12345)

	req = NewRequestFromMessage(msg)
	assert.NotNil(t, req)

	assert.Equal(t, Get, req.GetMessage().GetCode())
	// &net.UDPConn{}, &net.UDPAddr{}
	assert.NotNil(t, NewClientRequestFromMessage(msg, make(map[string]string), nil))

	req = NewRequest(MessageConfirmable, Get, 12345)
	assert.Equal(t, uint8(0), req.GetMessage().GetMessageType())

	req.SetConfirmable(false)
	assert.Equal(t, uint8(1), req.GetMessage().GetMessageType())
}
