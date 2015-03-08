package goap

import (
	"fmt"
	"log"
	"net"
	"time"
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

func (c *Client) Dial(nwNet string, host string) {
	udpAddr, err := net.ResolveUDPAddr(nwNet, host)
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

func (c *Client) Send(msg *Message) error {
	err := c.validate()

	if err != nil {
		return err
	}

	// Send message
	b := MessageToBytes(msg)
	_, err = c.conn.Write(b)

	if err != nil {
		log.Println(err)
	}

	if msg.MessageType == TYPE_NONCONFIRMABLE {
		return nil
	} else {
		// Read response
		var buf []byte = make([]byte, 1500)

		c.conn.SetReadDeadline(time.Now().Add(time.Second * 2))
		n, _, err := c.conn.ReadFromUDP(buf)

		if err != nil {
			return err
		}

		if c.successHandler != nil {
			resp, err := BytesToMessage(buf[:n])

			if err != nil {
				return err
			}

			fn := c.successHandler
			fn(resp)
		}
	}
	return nil
}

func (c *Client) Close() {
	defer c.conn.Close()
}
