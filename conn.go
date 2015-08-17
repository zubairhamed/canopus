package canopus

import (
	"errors"
	"net"
	"time"
)

// A simple wrapper interface around a connection
// This was primarily concieved so that mocks could be
// created to unit test connection code
type CanopusConnection interface {
	GetConnection() net.Conn
	Write(b []byte) (int, error)
	SetReadDeadline(t time.Time) error
	Read() (buf []byte, n int, err error)
	WriteTo(b []byte, addr net.Addr) (int, error)
}

// ----------------------------------------------------------------

func NewCanopusUDPConnection(c *net.UDPConn) CanopusConnection {
	return &CanopusUDPConnection{
		conn: c,
	}
}

func NewCanopusUDPConnectionWithAddr(c *net.UDPConn, a net.Addr) CanopusConnection {
	return &CanopusUDPConnection{
		conn: c,
		addr: a,
	}
}

type CanopusUDPConnection struct {
	conn net.Conn
	addr net.Addr
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
	n, _, err = c.conn.(*net.UDPConn).ReadFromUDP(buf)

	return
}

func (c *CanopusUDPConnection) WriteTo(b []byte, addr net.Addr) (n int, err error) {
	n, err = c.conn.(*net.UDPConn).WriteToUDP(b, addr.(*net.UDPAddr))

	return
}

// ----------------------------------------------------------------
func NewMockCanopusUDPConnection(code CoapCode, writeErr bool, readErr bool) *MockCanopusUDPConnection {
	return &MockCanopusUDPConnection{
		coapCode: code,
		writeErr: writeErr,
		readErr:  readErr,
	}
}

type MockCanopusUDPConnection struct {
	coapCode CoapCode
	writeErr bool
	readErr  bool
}

func (c *MockCanopusUDPConnection) GetConnection() net.Conn {
	return nil
}

func (c *MockCanopusUDPConnection) Write(b []byte) (n int, err error) {
	if c.writeErr {
		err = errors.New("Mock Write Error Generated")
	}
	n = len(b)
	return
}

func (c *MockCanopusUDPConnection) SetReadDeadline(t time.Time) error {
	return nil
}

func (c *MockCanopusUDPConnection) Read() (buf []byte, n int, err error) {
	msg := NewMessage(TYPE_NONCONFIRMABLE, c.coapCode, 12345)
	buf, _ = MessageToBytes(msg)
	n = len(buf)

	if c.readErr {
		err = errors.New("Mock Read Error Generated")
	}

	return
}

func (c *MockCanopusUDPConnection) WriteTo(b []byte, addr net.Addr) (n int, err error) {
	n = len(b)

	if c.writeErr {
		err = errors.New("Mock Write Error Generated")
	}
	return
}
