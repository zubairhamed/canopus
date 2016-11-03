package canopus

import "net"

type Session interface {
	GetMessage() (*Message, error)
	GetAddress() net.Addr
	Write([]byte)
	Read([]byte)
}

type UDPSession struct {
	addr net.Addr
	msgBuf []byte
}

type DTLSSession struct {

}
