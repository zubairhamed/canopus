package canopus

import (
	"net"
	"strconv"
	"strings"
)

// Creates a New Request Instance
func NewRequest(messageType uint8, messageMethod CoapCode, messageId uint16) CoapRequest {
	msg := NewMessage(messageType, messageMethod, messageId)
	msg.Token = []byte(GenerateToken(8))

	return &DefaultCoapRequest{
		msg: msg,
	}
}

func NewConfirmableGetRequest() CoapRequest {
	msg := NewMessage(MessageConfirmable, Get, GenerateMessageID())
	msg.Token = []byte(GenerateToken(8))

	return &DefaultCoapRequest{
		msg: msg,
	}
}

func NewConfirmablePostRequest() CoapRequest {
	msg := NewMessage(MessageConfirmable, Post, GenerateMessageID())
	msg.Token = []byte(GenerateToken(8))

	return &DefaultCoapRequest{
		msg: msg,
	}
}

func NewConfirmablePutRequest() CoapRequest {
	msg := NewMessage(MessageConfirmable, Put, GenerateMessageID())
	msg.Token = []byte(GenerateToken(8))

	return &DefaultCoapRequest{
		msg: msg,
	}
}

func NewConfirmableDeleteRequest() CoapRequest {
	msg := NewMessage(MessageConfirmable, Delete, GenerateMessageID())
	msg.Token = []byte(GenerateToken(8))

	return &DefaultCoapRequest{
		msg: msg,
	}
}

// Creates a new request messages from a CoAP Message
func NewRequestFromMessage(msg *Message) CoapRequest {
	return &DefaultCoapRequest{
		msg: msg,
	}
}

func NewClientRequestFromMessage(msg *Message, attrs map[string]string, conn *net.UDPConn, addr *net.UDPAddr) CoapRequest {
	return &DefaultCoapRequest{
		msg:   msg,
		attrs: attrs,
		conn:  conn,
		addr:  addr,
	}
}

type CoapRequest interface {
	SetProxyUri(uri string)
	SetMediaType(mt MediaType)
	GetConnection() *net.UDPConn
	GetAddress() *net.UDPAddr
	GetAttributes() map[string]string
	GetAttribute(o string) string
	GetAttributeAsInt(o string) int
	GetMessage() *Message
	SetStringPayload(s string)
	SetRequestURI(uri string)
	SetConfirmable(con bool)
	SetToken(t string)
	GetUriQuery(q string) string
	SetUriQuery(k string, v string)
}

// Wraps a CoAP Message as a Request
// Provides various methods which proxies the Message object methods
type DefaultCoapRequest struct {
	msg    *Message
	attrs  map[string]string
	conn   *net.UDPConn
	addr   *net.UDPAddr
	server *CoapServer
}

func (c *DefaultCoapRequest) SetProxyUri(uri string) {
	c.msg.AddOption(OptionProxyURI, uri)
}

func (c *DefaultCoapRequest) SetMediaType(mt MediaType) {
	c.msg.AddOption(OptionContentFormat, mt)
}

func (c *DefaultCoapRequest) GetConnection() *net.UDPConn {
	return c.conn
}

func (c *DefaultCoapRequest) GetAddress() *net.UDPAddr {
	return c.addr
}

func (c *DefaultCoapRequest) GetAttributes() map[string]string {
	return c.attrs
}

func (c *DefaultCoapRequest) GetAttribute(o string) string {
	return c.attrs[o]
}

func (c *DefaultCoapRequest) GetAttributeAsInt(o string) int {
	attr := c.GetAttribute(o)
	i, _ := strconv.Atoi(attr)

	return i
}

func (c *DefaultCoapRequest) GetMessage() *Message {
	return c.msg
}

func (c *DefaultCoapRequest) SetStringPayload(s string) {
	c.msg.Payload = NewPlainTextPayload(s)
}

func (c *DefaultCoapRequest) SetRequestURI(uri string) {
	c.msg.AddOptions(NewPathOptions(uri))
}

func (c *DefaultCoapRequest) SetConfirmable(con bool) {
	if con {
		c.msg.MessageType = MessageConfirmable
	} else {
		c.msg.MessageType = MessageNonConfirmable
	}
}

func (c *DefaultCoapRequest) SetToken(t string) {
	c.msg.Token = []byte(t)
}

func (c *DefaultCoapRequest) GetUriQuery(q string) string {
	qs := c.GetMessage().GetOptionsAsString(OptionURIQuery)

	for _, o := range qs {
		ps := strings.Split(o, "=")
		if len(ps) == 2 {
			if ps[0] == q {
				return ps[1]
			}
		}
	}
	return ""
}

func (c *DefaultCoapRequest) SetUriQuery(k string, v string) {
	c.GetMessage().AddOption(OptionURIQuery, k+"="+v)
}
