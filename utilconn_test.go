package canopus

import (
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func TestSendMessages(t *testing.T) {
	var conn CanopusConnection
	_, err := SendMessageTo(nil, conn, nil)
	assert.NotNil(t, err)
	assert.Equal(t, ErrNilConn, err)

	conn = NewCanopusUDPConnection(nil)
	SendMessageTo(nil, conn, nil)
	_, err = SendMessageTo(nil, conn, nil)
	assert.NotNil(t, err)
	assert.Equal(t, ErrNilMessage, err)

	_, err = SendMessageTo(NewEmptyMessage(12345), conn, nil)
	assert.NotNil(t, err)
	assert.Equal(t, ErrNilAddr, err)

	addr := &net.UDPAddr{}
	conn = NewMockCanopusUDPConnection(CoapCode_Created, false, false)
	msg := NewBasicConfirmableMessage()
	_, err = SendMessageTo(msg, conn, addr)
	assert.Nil(t, err)

	msg.MessageType = MessageNonConfirmable
	_, err = SendMessageTo(msg, conn, addr)
	assert.Nil(t, err)

	conn = NewMockCanopusUDPConnection(CoapCode_Created, false, true)
	msg.MessageType = MessageConfirmable
	_, err = SendMessageTo(msg, conn, addr)
	assert.NotNil(t, err)
}
