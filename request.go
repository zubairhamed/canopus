package canopus

import (
	"strconv"
	"strings"
)

// Creates a New Request Instance
func NewRequest(messageType uint8, messageMethod CoapCode, messageID uint16) Request {
	msg := NewMessage(messageType, messageMethod, messageID)
	msg.Token = []byte(GenerateToken(8))

	return &DefaultCoapRequest{
		msg: msg,
	}
}

func NewConfirmableGetRequest() Request {
	msg := NewMessage(MessageConfirmable, Get, GenerateMessageID())
	msg.Token = []byte(GenerateToken(8))

	return &DefaultCoapRequest{
		msg: msg,
	}
}

func NewConfirmablePostRequest() Request {
	msg := NewMessage(MessageConfirmable, Post, GenerateMessageID())
	msg.Token = []byte(GenerateToken(8))

	return &DefaultCoapRequest{
		msg: msg,
	}
}

func NewConfirmablePutRequest() Request {
	msg := NewMessage(MessageConfirmable, Put, GenerateMessageID())
	msg.Token = []byte(GenerateToken(8))

	return &DefaultCoapRequest{
		msg: msg,
	}
}

func NewConfirmableDeleteRequest() Request {
	msg := NewMessage(MessageConfirmable, Delete, GenerateMessageID())
	msg.Token = []byte(GenerateToken(8))

	return &DefaultCoapRequest{
		msg: msg,
	}
}

// Creates a new request messages from a CoAP Message
func NewRequestFromMessage(msg Message) Request {
	return &DefaultCoapRequest{
		msg: msg,
	}
}

func NewClientRequestFromMessage(msg Message, attrs map[string]string, session Session) Request {
	return &DefaultCoapRequest{
		msg:     msg,
		attrs:   attrs,
		session: session,
	}
}

// Wraps a CoAP Message as a Request
// Provides various methods which proxies the Message object methods
type DefaultCoapRequest struct {
	msg     Message
	attrs   map[string]string
	session Session
	server  *CoapServer
}

func (c *DefaultCoapRequest) SetProxyURI(uri string) {
	c.msg.AddOption(OptionProxyURI, uri)
}

func (c *DefaultCoapRequest) SetMediaType(mt MediaType) {
	c.msg.AddOption(OptionContentFormat, mt)
}

func (c *DefaultCoapRequest) GetSession() Session {
	return c.session
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

func (c *DefaultCoapRequest) GetMessage() Message {
	return c.msg
}

func (c *DefaultCoapRequest) SetStringPayload(s string) {
	c.msg.SetPayload(payload.NewPlainTextPayload(s))
}

func (c *DefaultCoapRequest) SetPayload(b []byte) {
	c.msg.SetPayload(payload.NewBytesPayload(b))
}

func (c *DefaultCoapRequest) SetRequestURI(uri string) {
	c.msg.AddOptions(NewPathOptions(uri))
}

func (c *DefaultCoapRequest) SetConfirmable(con bool) {
	if con {
		c.msg.GetMessageType() = MessageConfirmable
	} else {
		c.msg.GetMessageType() = MessageNonConfirmable
	}
}

func (c *DefaultCoapRequest) SetToken(t string) {
	c.msg.SetToken([]byte(t))
}

func (c *DefaultCoapRequest) GetURIQuery(q string) string {
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

func (c *DefaultCoapRequest) SetURIQuery(k string, v string) {
	c.GetMessage().AddOption(OptionURIQuery, k+"="+v)
}
