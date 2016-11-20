package canopus

import (
	"fmt"
	"net"
)

type UDPServerSession struct {
	addr   net.Addr
	conn   ServerConnection
	server CoapServer
	rcvd   chan []byte
}

func (s *UDPServerSession) GetConnection() ServerConnection {
	return s.conn
}

func (s *UDPServerSession) GetAddress() net.Addr {
	return s.addr
}

func (s *UDPServerSession) Received(b []byte) (n int) {
	l := len(b)
	go func() {
		s.rcvd <- b
	}()
	return l
}

func (s *UDPServerSession) Write(b []byte) (n int, err error) {
	n, err = s.conn.WriteTo(b, s.GetAddress())

	return
}

func (s *UDPServerSession) Read(b []byte) (n int, err error) {
	fmt.Println("UDPServerSession:Read")
	data := <-s.rcvd
	fmt.Println("UDPServerSession:Read2")
	copy(b, data)
	fmt.Println("UDPServerSession:Read3")
	return len(data), nil
}

func (s *UDPServerSession) GetServer() CoapServer {
	return s.server
}
