package canopus

import (
	"net"
	"time"
)

// Sends a CoAP Message to UDP address
func SendMessageTo(msg *Message, conn CanopusConnection, addr *net.UDPAddr) (CoapResponse, error) {
	if conn == nil {
		return nil, ERR_NIL_CONN
	}

	if msg == nil {
		return nil, ERR_NIL_MESSAGE
	}

	if addr == nil {
		return nil, ERR_NIL_ADDR
	}

	b, _ := MessageToBytes(msg)
	_, err := conn.WriteTo(b, addr)

	if err != nil {
		return nil, err
	}

	if msg.MessageType == TYPE_NONCONFIRMABLE {
		return NewResponse(NewEmptyMessage(msg.MessageId), err), err
	} else {
		// conn.SetReadDeadline(time.Now().Add(time.Second * 2))
		buf, n, err := conn.Read()
		if err != nil {
			return nil, err
		}
		msg, err := BytesToMessage(buf[:n])
		resp := NewResponse(msg, err)

		return resp, err
	}
	return nil, nil
}

// Sends a CoAP Message to a UDP Connection
func SendMessage(msg *Message, conn CanopusConnection) (CoapResponse, error) {
	if conn == nil {
		return nil, ERR_NIL_CONN
	}

	b, _ := MessageToBytes(msg)
	_, err := conn.Write(b)

	if err != nil {
		return nil, err
	}

	if msg.MessageType == TYPE_NONCONFIRMABLE {
		return nil, err
	} else {
		var buf []byte = make([]byte, 1500)
		conn.SetReadDeadline(time.Now().Add(time.Second * 2))
		buf, n, err := conn.Read()

		if err != nil {
			return nil, err
		}

		msg, err := BytesToMessage(buf[:n])

		resp := NewResponse(msg, err)

		return resp, err
	}
}
