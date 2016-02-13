package canopus

import (
	"net"
	"time"
)

// SendMessageTo sends a CoAP Message to UDP address
func SendMessageTo(msg *Message, conn Connection, addr *net.UDPAddr) (CoapResponse, error) {
	if conn == nil {
		return nil, ErrNilConn
	}

	if msg == nil {
		return nil, ErrNilMessage
	}

	if addr == nil {
		return nil, ErrNilAddr
	}

	b, _ := MessageToBytes(msg)
	_, err := conn.WriteTo(b, addr)

	if err != nil {
		return nil, err
	}

	if msg.MessageType == MessageNonConfirmable {
		return NewResponse(NewEmptyMessage(msg.MessageID), err), err
	}

	// conn.SetReadDeadline(time.Now().Add(time.Second * 2))
	buf, n, err := conn.Read()
	if err != nil {
		return nil, err
	}
	msg, err = BytesToMessage(buf[:n])
	resp := NewResponse(msg, err)

	return resp, err
}

// SendMessage sends a CoAP Message to a UDP Connection
func SendMessage(msg *Message, conn Connection) (CoapResponse, error) {
	if conn == nil {
		return nil, ErrNilConn
	}

	b, _ := MessageToBytes(msg)
	_, err := conn.Write(b)

	if err != nil {
		return nil, err
	}

	if msg.MessageType == MessageNonConfirmable {
		return nil, err
	}

	var buf = make([]byte, 1500)
	conn.SetReadDeadline(time.Now().Add(time.Second * DefaultAckTimeout))
	buf, n, err := conn.Read()

	if err != nil {
		return nil, err
	}

	msg, err = BytesToMessage(buf[:n])

	resp := NewResponse(msg, err)

	return resp, err
}
