package canopus

import (
	"strconv"
	"strings"
)

// Creates a New Request Instance
func NewRequest(messageType uint8, messageMethod CoapCode) Request {
	return NewRequestWithMessageId(messageType, messageMethod, GenerateMessageID())
}

func NewRequestWithMessageId(messageType uint8, messageMethod CoapCode, messageID uint16) Request {
	msg := NewMessage(messageType, messageMethod, messageID)
	return &CoapRequest{
		msg: msg,
	}
}

func NewConfirmableGetRequest() Request {
	return &CoapRequest{
		msg: NewMessage(MessageConfirmable, Get, GenerateMessageID()),
	}
}

func NewConfirmablePostRequest() Request {
	return &CoapRequest{
		msg: NewMessage(MessageConfirmable, Post, GenerateMessageID()),
	}
}

func NewConfirmablePutRequest() Request {
	return &CoapRequest{
		msg: NewMessage(MessageConfirmable, Put, GenerateMessageID()),
	}
}

func NewConfirmableDeleteRequest() Request {
	return &CoapRequest{
		msg: NewMessage(MessageConfirmable, Delete, GenerateMessageID()),
	}
}

// Creates a new request messages from a CoAP Message
func NewRequestFromMessage(msg Message) Request {
	return &CoapRequest{
		msg: msg,
	}
}

func NewClientRequestFromMessage(msg Message, attrs map[string]string, session Session) Request {
	return &CoapRequest{
		msg:     msg,
		attrs:   attrs,
		session: session,
	}
}

// Wraps a CoAP Message as a Request
// Provides various methods which proxies the Message object methods
type CoapRequest struct {
	msg     Message
	attrs   map[string]string
	session Session
	server  *CoapServer
}

func (c *CoapRequest) SetProxyURI(uri string) {
	c.msg.AddOption(OptionProxyURI, uri)
}

func (c *CoapRequest) SetMediaType(mt MediaType) {
	c.msg.AddOption(OptionContentFormat, mt)
}

func (c *CoapRequest) GetSession() Session {
	return c.session
}

func (c *CoapRequest) GetAttributes() map[string]string {
	return c.attrs
}

func (c *CoapRequest) GetAttribute(o string) string {
	return c.attrs[o]
}

func (c *CoapRequest) GetAttributeAsInt(o string) int {
	attr := c.GetAttribute(o)
	i, _ := strconv.Atoi(attr)

	return i
}

func (c *CoapRequest) GetMessage() Message {
	return c.msg
}

func (c *CoapRequest) SetStringPayload(s string) {
	c.msg.(*CoapMessage).SetPayload(NewPlainTextPayload(s))
}

func (c *CoapRequest) SetPayload(b []byte) {
	c.msg.(*CoapMessage).SetPayload(NewBytesPayload(b))
}

func (c *CoapRequest) SetRequestURI(uri string) {
	c.msg.AddOptions(NewPathOptions(uri))
}

func (c *CoapRequest) SetConfirmable(con bool) {
	if con {
		c.msg.(*CoapMessage).SetMessageType(MessageConfirmable)
	} else {
		c.msg.(*CoapMessage).SetMessageType(MessageNonConfirmable)
	}
}

func (c *CoapRequest) SetToken(t string) {
	c.msg.(*CoapMessage).SetToken([]byte(t))
}

func (c *CoapRequest) GetURIQuery(q string) string {
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

func (c *CoapRequest) SetURIQuery(k string, v string) {
	c.GetMessage().AddOption(OptionURIQuery, k+"="+v)
}
