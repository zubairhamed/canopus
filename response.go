package canopus

import "strings"

// Creates a new Response object with a Message object and any error messages
func NewResponse(msg *Message, err error) *Response {
	resp := &Response{
		msg: msg,
		err: err,
	}

	return resp
}

// Creates a new response object with a Message object
func NewResponseWithMessage(msg *Message) *Response {
	resp := &Response{
		msg: msg,
	}

	return resp
}

type Response struct {
	msg *Message
	err error
}

func (c *Response) GetMessage() *Message {
	return c.msg
}

func (c *Response) GetError() error {
	return c.err
}

func (c *Response) GetPayload() []byte {
	return c.GetMessage().Payload.GetBytes()
}

func (c *Response) GetUriQuery(q string) string {
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
