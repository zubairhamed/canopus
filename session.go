package canopus

import "net"

type ServerSession struct {
	buf    []byte
	addr   net.Addr
	conn   CanopusConnection
	server CoapServer
}

func (s *ServerSession) GetConnection() CanopusConnection {
	return s.conn
}

func (s *ServerSession) GetAddress() net.Addr {
	return s.addr
}

func (s *ServerSession) Write(b []byte) {
	s.buf = append(s.buf, b...)
}

func (s *ServerSession) FlushBuffer() {
	s.buf = nil
}

func (s *ServerSession) Read() []byte {
	return s.buf
}

func (s *ServerSession) GetServer() CoapServer {
	return s.server
}
