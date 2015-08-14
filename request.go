package canopus

import (
	"net"
	"strconv"
	"strings"
)

// Creates a New Request Instance
func NewRequest(messageType uint8, messageMethod CoapCode, messageId uint16) *Request {
	msg := NewMessage(messageType, messageMethod, messageId)
	msg.Token = []byte(GenerateToken(8))

	return &Request{
		msg: msg,
	}
}

// Creates a new request messages from a CoAP Message
func NewRequestFromMessage(msg *Message) *Request {
	return &Request{
		msg: msg,
	}
}

func NewClientRequestFromMessage(msg *Message, attrs map[string]string, conn *net.UDPConn, addr *net.UDPAddr) *Request {
	return &Request{
		msg:   msg,
		attrs: attrs,
		conn:  conn,
		addr:  addr,
	}
}

// Wraps a CoAP Message as a Request
// Provides various methods which proxies the Message object methods
type Request struct {
	msg    *Message
	attrs  map[string]string
	conn   *net.UDPConn
	addr   *net.UDPAddr
	server *CoapServer
}

func (c *Request) SetProxyUri(uri string) {
	c.msg.AddOption(OPTION_PROXY_URI, uri)
}

func (c *Request) SetMediaType(mt MediaType) {
	c.msg.AddOption(OPTION_CONTENT_FORMAT, mt)
}

func (c *Request) GetConnection() *net.UDPConn {
	return c.conn
}

func (c *Request) GetAddress() *net.UDPAddr {
	return c.addr
}

func (c *Request) GetAttributes() map[string]string {
	return c.attrs
}

func (c *Request) GetAttribute(o string) string {
	return c.attrs[o]
}

func (c *Request) GetAttributeAsInt(o string) int {
	attr := c.GetAttribute(o)
	i, _ := strconv.Atoi(attr)

	return i
}

func (c *Request) GetMessage() *Message {
	return c.msg
}

func (c *Request) SetStringPayload(s string) {
	c.msg.Payload = NewPlainTextPayload(s)
}

func (c *Request) SetRequestURI(uri string) {
	c.msg.AddOptions(NewPathOptions(uri))
}

func (c *Request) SetConfirmable(con bool) {
	if con {
		c.msg.MessageType = TYPE_CONFIRMABLE
	} else {
		c.msg.MessageType = TYPE_NONCONFIRMABLE
	}
}

func (c *Request) SetToken(t string) {
	c.msg.Token = []byte(t)
}

func (c *Request) IncrementMessageId() {
	c.msg.MessageId = c.msg.MessageId + 1
}

func (c *Request) GetUriQuery(q string) string {
	qs := c.GetMessage().GetOptionsAsString(OPTION_URI_QUERY)

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

func (c *Request) SetUriQuery(k string, v string) {
	c.GetMessage().AddOption(OPTION_URI_QUERY, k+"="+v)
}

func (c *Request) GetServer() *CoapServer {
	return c.server
}
