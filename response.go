package canopus

import "strings"

func NoResponse() CoapResponse {
	return NilResponse{}
}

type CoapResponse interface {
	GetMessage() *Message
	GetError() error
	GetPayload() []byte
	GetURIQuery(q string) string
}

type NilResponse struct {
}

func (c NilResponse) GetMessage() *Message {
	return nil
}

func (c NilResponse) GetError() error {
	return nil
}

func (c NilResponse) GetPayload() []byte {
	return nil
}

func (c NilResponse) GetURIQuery(q string) string {
	return ""
}

// Creates a new Response object with a Message object and any error messages
func NewResponse(msg *Message, err error) CoapResponse {
	resp := &DefaultResponse{
		msg: msg,
		err: err,
	}

	return resp
}

// Creates a new response object with a Message object
func NewResponseWithMessage(msg *Message) CoapResponse {
	resp := &DefaultResponse{
		msg: msg,
	}

	return resp
}

type DefaultResponse struct {
	msg *Message
	err error
}

func (c *DefaultResponse) GetMessage() *Message {
	return c.msg
}

func (c *DefaultResponse) GetError() error {
	return c.err
}

func (c *DefaultResponse) GetPayload() []byte {
	return c.GetMessage().Payload.GetBytes()
}

func (c *DefaultResponse) GetURIQuery(q string) string {
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
