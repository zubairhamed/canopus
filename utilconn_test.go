package canopus

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"net"
)

func TestSendMessages(t *testing.T) {
	var conn CanopusConnection
	_, err := SendMessageTo(nil, conn, nil)
	assert.NotNil(t, err)
	assert.Equal(t, ERR_NIL_CONN, err)

	conn = NewCanopusUDPConnection(nil)
	SendMessageTo(nil, conn, nil)
	_, err = SendMessageTo(nil, conn, nil)
	assert.NotNil(t, err)
	assert.Equal(t, ERR_NIL_MESSAGE, err)

	_, err = SendMessageTo(NewEmptyMessage(12345), conn, nil)
	assert.NotNil(t, err)
	assert.Equal(t, ERR_NIL_ADDR, err)

	addr := &net.UDPAddr{}
	conn = NewMockCanopusUDPConnection(COAPCODE_201_CREATED)
	msg := NewBasicConfirmableMessage()
	_, err = SendMessageTo(msg, conn, addr)
	assert.Nil(t, err)
}
