package canopus

import "net"

type Session struct {
	buf    []byte
	addr   net.Addr
	conn   CanopusConnection
	server CoapServer
}

func (s *Session) GetConnection() CanopusConnection {
	return s.conn
}

func (s *Session) GetAddress() net.Addr {
	return s.addr
}

func (s *Session) Write(b []byte) {
	s.buf = append(s.buf, b...)
}

func (s *Session) FlushBuffer() {
	s.buf = nil
}

func (s *Session) Read() []byte {
	return s.buf
}

func (s *Session) GetServer() CoapServer {
	return s.server
}
