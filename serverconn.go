package canopus

import (
	"net"
	"time"
)

type DTLSServerConnection struct {
}

type UDPServerConnection struct {
	conn net.PacketConn
}

func (uc *UDPServerConnection) ReadFrom(b []byte) (n int, addr net.Addr, err error) {
	return uc.conn.ReadFrom(b)
}

func (uc *UDPServerConnection) WriteTo(b []byte, addr net.Addr) (n int, err error) {
	return uc.conn.WriteTo(b, addr)
}

func (uc *UDPServerConnection) Close() error {
	return uc.conn.Close()
}

func (uc *UDPServerConnection) LocalAddr() net.Addr {
	return uc.conn.LocalAddr()
}

func (uc *UDPServerConnection) SetDeadline(t time.Time) error {
	return uc.conn.SetDeadline(t)
}

func (uc *UDPServerConnection) SetReadDeadline(t time.Time) error {
	return uc.conn.SetReadDeadline(t)
}

func (uc *UDPServerConnection) SetWriteDeadline(t time.Time) error {
	return uc.conn.SetWriteDeadline(t)
}
