package goap

import (
	"fmt"
	"log"
	"net"
	"strconv"
)

func NewClient() *Client {
	return &Client{
		async: true,
	}
}

type Client struct {
	conn           *net.UDPConn
	async          bool
	successHandler MessageHandler
	timeoutHandler MessageHandler
	resetHandler   MessageHandler
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

func (c *Client) OnSuccess(fn MessageHandler) {
	fmt.Println()
	c.successHandler = fn
}

func (c *Client) OnReset(fn MessageHandler) {
	c.resetHandler = fn
}

func (c *Client) OnTimeout(fn MessageHandler) {
	c.timeoutHandler = fn
}

func (c *Client) validate() error {
	return nil
}

func (c *Client) doSend(msg *Message) (*Message, error) {
    resp, err := SendMessage(msg, c.conn)

    return resp, err
}

func (c *Client) Send(msg *Message)(*Message, error) {
    return c.doSend (msg)
}

func (c *Client) SendAsync(msg *Message, fn MessageHandler) {
    msg, err := c.doSend(msg)

    fn (msg, err)
}

func (c *Client) Close() {
	defer c.conn.Close()
}
