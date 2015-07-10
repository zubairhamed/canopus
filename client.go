package canopus

import (
	"github.com/zubairhamed/go-commons/logging"
	"net"
)

func NewCoapClient(local string) *CoapClient {
	localAddr, err := net.ResolveUDPAddr("udp", local)
	if err != nil {
		logging.LogError("Error starting CoAP Server: ", err)
	}

	return &CoapClient{
		localAddr: localAddr,
	}
}

type CoapClient struct {
	localAddr  *net.UDPAddr
	remoteAddr *net.UDPAddr
	conn       *net.UDPConn
}

func (c *CoapClient) Dial(host string) {
	remAddr, err := net.ResolveUDPAddr("udp", host)
	logging.LogError(err)

	c.remoteAddr = remAddr

	conn, err := net.DialUDP("udp", c.localAddr, c.remoteAddr)
	logging.LogError(err)

	c.conn = conn
}

func (c *CoapClient) doSend(req *CoapRequest, conn *net.UDPConn, fn ResponseHandler) (resp *CoapResponse, err error) {
	if fn == nil {
		resp, err = SendMessage(req.GetMessage(), conn)
		return
	} else {
		SendAsyncMessage(req.GetMessage(), conn, fn)
		return
	}
}

func (c *CoapClient) Send(req *CoapRequest) (*CoapResponse, error) {
	return c.doSend(req, c.conn, nil)
}

func (c *CoapClient) SendTo(req *CoapRequest, conn *net.UDPConn) (*CoapResponse, error) {
	return c.doSend(req, conn, nil)
}

func (c *CoapClient) SendAsync(req *CoapRequest, fn ResponseHandler) {
	c.doSend(req, c.conn, fn)
}

func (c *CoapClient) Close() {
	c.conn.Close()
}
