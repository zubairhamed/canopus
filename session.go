package canopus

import "net"

type UDPServerSession struct {
	addr   net.Addr
	conn   ServerConnection
	server CoapServer
	rcvd   chan []byte
	buf    []byte
}

func (s *UDPServerSession) GetConnection() ServerConnection {
	return s.conn
}

func (s *UDPServerSession) GetAddress() net.Addr {
	return s.addr
}

func (s *UDPServerSession) WriteBuffer(b []byte) (n int) {
	l := len(b)
	s.buf = append(s.buf, b...)
	return l
}

func (s *UDPServerSession) Write(b []byte) (n int, err error) {
	n, err = s.conn.WriteTo(b, s.GetAddress())

	return
}

func (s *UDPServerSession) Read(b []byte) (n int, err error) {
	data := <-s.rcvd
	copy(b, data)
	return len(data), nil
}

func (s *UDPServerSession) GetServer() CoapServer {
	return s.server
}
