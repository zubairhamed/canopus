package goap

import (
	"log"
	"net"
	"strconv"
)

func NewClient() *Client {
	return &Client{}
}

type Client struct {
	conn           *net.UDPConn
}

func (c *Client) Dial(nwNet string, host string, port int) {
	hostString := host + ":" + strconv.Itoa(port)
	udpAddr, err := net.ResolveUDPAddr(nwNet, hostString)
	if err != nil {
		log.Println(err)
	}

	conn, err := net.DialUDP(nwNet, nil, udpAddr)
	if err != nil {
		log.Println(err)
	}
	c.conn = conn
}

func (c *Client) doSend(req *CoapRequest) (*CoapResponse, error) {
    resp, err := SendMessage(req.GetMessage(), c.conn)

    return resp, err
}

func (c *Client) Send(req *CoapRequest)(*CoapResponse, error) {
    return c.doSend (req)
}

func (c *Client) SendAsync(req *CoapRequest, fn ResponseHandler) {
    resp, err := c.doSend(req)

    fn (resp, err)
}

/*
func (c *Client) Discover(fn ResponseHandler) {
    // TODO: Construct Discovery Payload
    req := nil
    resp, err := c.doSend(req)

    fn (resp, err)
}
*/

func (c *Client) Close() {
	defer c.conn.Close()
}
