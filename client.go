package canopus

import (
	"log"
	"net"
	"github.com/zubairhamed/go-commons/logging"
)

type CoapClient struct {
	localAddr  *net.UDPAddr
	remoteAddr *net.UDPAddr
	conn       *net.UDPConn
}

func (c CoapClient) Dial(host string) {
	remAddr, err := net.ResolveUDPAddr("udp", host)
	logging.LogError(err)

	c.remoteAddr = remAddr

	conn, err := net.DialUDP("udp", c.localAddr, c.remoteAddr)
	logging.LogError(err)

	c.conn = conn
}

func (c *CoapClient) doSend(req *CoapRequest, conn *net.UDPConn) (*CoapResponse, error) {
	log.Println(req, conn)
	resp, err := SendMessage(req.GetMessage(), conn)

	return resp, err
}

func (c *CoapClient) Send(req *CoapRequest) (*CoapResponse, error) {
	log.Println("@@", req, c.conn)
	return c.doSend(req, c.conn)
}

func (c *CoapClient) SendTo(req *CoapRequest, conn *net.UDPConn) (*CoapResponse, error) {
	return c.doSend(req, conn)
}

func (c *CoapClient) SendAsync(req *CoapRequest, fn ResponseHandler) {

}

func (c *CoapClient) Close() {
	c.conn.Close()
}
