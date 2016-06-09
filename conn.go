package canopus

import (
	"errors"
	"net"
	"time"
)

// ----------------------------------------------------------------

// NewCanopusUDPConnection creates a new default CanopousConnect
func NewUDPConnection(c *net.UDPConn) Connection {
	return &UDPConnection{
		conn: c,
	}
}

type UDPConnection struct {
	conn *net.UDPConn
	addr net.Addr
}

func (c *UDPConnection) GetConnection() net.Conn {
	return c.conn
}

func (c *UDPConnection) Write(b []byte) (int, error) {
	return c.conn.Write(b)
}

func (c *UDPConnection) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

func (c *UDPConnection) Read() (buf []byte, n int, err error) {
	buf = make([]byte, MaxPacketSize)
	n, _, err = c.conn.ReadFromUDP(buf)

	return
}

func (c *UDPConnection) WriteTo(b []byte, addr net.Addr) (n int, err error) {
	n, err = c.conn.WriteToUDP(b, addr.(*net.UDPAddr))

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
	msg := NewMessage(MessageNonConfirmable, c.coapCode, 12345)
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
