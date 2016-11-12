package canopus

import "net"

type UDPServerSession struct {
	buf    []byte
	addr   net.Addr
	conn   ServerConnection
	server CoapServer
}

func (s *UDPServerSession) GetConnection() ServerConnection {
	return s.conn
}

func (s *UDPServerSession) GetAddress() net.Addr {
	return s.addr
}

func (s *UDPServerSession) Received(b []byte) (n int) {
	n = len(b)
	s.buf = append(s.buf, b...)

	return
}

func (s *UDPServerSession) Write(b []byte) (n int, err error) {
	n, err = s.conn.WriteTo(b, s.GetAddress())
	s.flushBuffer()

	return
}

func (s *UDPServerSession) flushBuffer() {
	s.buf = nil
}

func (s *UDPServerSession) Read(b []byte) (n int, err error) {
	b = s.buf

	return len(b), nil
}

func (s *UDPServerSession) GetServer() CoapServer {
	return s.server
}
