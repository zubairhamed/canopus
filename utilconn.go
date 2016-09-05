package canopus

import (
	"net"
)

type CoapResponseChannel struct {
	Response CoapResponse
	Error		 error
}

func doSendMessage(c CoapServer, msg *Message, conn Connection, addr *net.UDPAddr, ch chan *CoapResponseChannel) {
	resp := &CoapResponseChannel{}

	b, err := MessageToBytes(msg)
	if err != nil {
		resp.Error = err
		ch <- resp
	}

	_, err = conn.WriteTo(b, addr)
	if err != nil {
		resp.Error = err
		ch <- resp
	}

	if msg.MessageType == MessageNonConfirmable {
		resp.Response = NewResponse(NewEmptyMessage(msg.MessageID), nil)
		ch <- resp
	}

	AddResponseChannel(c, msg.MessageID, ch)
}

// SendMessageTo sends a CoAP Message to UDP address
func SendMessageTo(c CoapServer, msg *Message, conn Connection, addr *net.UDPAddr) (CoapResponse, error) {
	if conn == nil {
		return nil, ErrNilConn
	}

	if msg == nil {
		return nil, ErrNilMessage
	}

	if addr == nil {
		return nil, ErrNilAddr
	}

	ch := NewResponseChannel()
	go doSendMessage(c, msg, conn, addr, ch)
	respCh := <- ch

	return respCh.Response, respCh.Error
}

func MessageSizeAllowed(req CoapRequest) bool {
	msg := req.GetMessage()
	b, _ := MessageToBytes(msg)

	if len(b) > 65536 {
		return false
	}

	return true
}
