package canopus

import (
	"net"
	"time"

	"github.com/jvermillard/nativedtls"
)

type DTLSServerConnection struct {
	ctx  *nativedtls.DTLSCtx
	conn net.PacketConn
}

func (uc *DTLSServerConnection) ReadFrom(b []byte) (n int, addr net.Addr, err error) {
	return 0, nil, nil
}

func (uc *DTLSServerConnection) WriteTo(b []byte, addr net.Addr) (n int, err error) {
	return 0, nil
}

func (uc *DTLSServerConnection) Close() error {
	return nil
}

func (uc *DTLSServerConnection) LocalAddr() net.Addr {
	return nil
}

func (uc *DTLSServerConnection) SetDeadline(t time.Time) error {
	return nil
}

func (uc *DTLSServerConnection) SetReadDeadline(t time.Time) error {
	return nil
}

func (uc *DTLSServerConnection) SetWriteDeadline(t time.Time) error {
	return nil
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
