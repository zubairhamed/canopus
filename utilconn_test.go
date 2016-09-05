package canopus

import (
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func TestSendMessages(t *testing.T) {
	var conn Connection
	var s CoapServer
	_, err := SendMessageTo(s, nil, conn, nil)
	assert.NotNil(t, err)
	assert.Equal(t, ErrNilConn, err)

	conn = NewUDPConnection(nil)
	SendMessageTo(s, nil, conn, nil)
	_, err = SendMessageTo(s, nil, conn, nil)
	assert.NotNil(t, err)
	assert.Equal(t, ErrNilMessage, err)

	_, err = SendMessageTo(s, NewEmptyMessage(12345), conn, nil)
	assert.NotNil(t, err)
	assert.Equal(t, ErrNilAddr, err)

	addr := &net.UDPAddr{}
	conn = NewMockCanopusUDPConnection(CoapCodeCreated, false, false)
	msg := NewBasicConfirmableMessage()
	_, err = SendMessageTo(s, msg, conn, addr)
	assert.Nil(t, err)

	msg.MessageType = MessageNonConfirmable
	_, err = SendMessageTo(s, msg, conn, addr)
	assert.Nil(t, err)

	conn = NewMockCanopusUDPConnection(CoapCodeCreated, false, true)
	msg.MessageType = MessageConfirmable
	_, err = SendMessageTo(s, msg, conn, addr)
	assert.NotNil(t, err)
}
