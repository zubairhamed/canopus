package canopus

import (
	"net"
	"time"
)

type CanopusConnection interface {
	GetConnection() net.Conn
	Write(b []byte) (int, error)
	SetReadDeadline(t time.Time) error
	Read() (buf []byte, n int, err error)
	WriteTo(b []byte, addr net.Addr) (int, error)

	// Send(*Message)(*Response, error)
	// SendTo(*Message, net.Addr)(*Response, error)
}

// ----------------------------------------------------------------

func NewCanopusUDPConnection(c *net.UDPConn) *CanopusUDPConnection {
	return &CanopusUDPConnection{
		conn: c,
	}
}

func NewCanopusUDPConnectionWithAddr(c *net.UDPConn, a net.Addr) *CanopusUDPConnection {
	return &CanopusUDPConnection{
		conn: c,
		addr: a,
	}
}

type CanopusUDPConnection struct {
	conn 	net.Conn
	addr 	net.Addr
}

func (c *CanopusUDPConnection) GetConnection() net.Conn {
	return c.conn
}

func (c *CanopusUDPConnection) Write(b []byte) (int, error) {
	return c.conn.(*net.UDPConn).Write(b)
}

func (c *CanopusUDPConnection) SetReadDeadline(t time.Time) error {
	return c.conn.(*net.UDPConn).SetReadDeadline(t)
}

func (c *CanopusUDPConnection) Read() (buf []byte, n int, err error) {
	buf = make([]byte, 1500)
	n, _,err = c.conn.(*net.UDPConn).ReadFromUDP(buf)

	return
}

func (c *CanopusUDPConnection) WriteTo(b []byte, addr net.Addr) (n int, err error) {
	n,err = c.conn.(*net.UDPConn).WriteToUDP(b, addr.(*net.UDPAddr))

	return
}

// ----------------------------------------------------------------
func NewMockCanopusUDPConnection(code CoapCode) *MockCanopusUDPConnection {
	return &MockCanopusUDPConnection{
		coapCode: code,
	}
}


type MockCanopusUDPConnection struct {
	coapCode 	CoapCode
}

func (c *MockCanopusUDPConnection) GetConnection() net.Conn {
	return nil
}

func (c *MockCanopusUDPConnection) Write(b []byte) (int, error) {
	return len(b), nil
}

func (c *MockCanopusUDPConnection) SetReadDeadline(t time.Time) error {
	return nil
}

func (c *MockCanopusUDPConnection) Read() (buf []byte, n int, err error) {
	msg := NewMessage(TYPE_NONCONFIRMABLE, c.coapCode, 12345)
	buf, _ = MessageToBytes(msg)
	n = len(buf)
	return
}

func (c *MockCanopusUDPConnection) WriteTo(b []byte, addr net.Addr) (n int, err error) {
	n = len(b)
	return
}

