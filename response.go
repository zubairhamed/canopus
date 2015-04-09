package goap

func NewResponse(msg *Message, err error) *CoapResponse {
    resp := &CoapResponse{
        msg: msg,
        err: err,
    }

    return resp
}

type CoapResponse struct {
    msg     *Message
    err     error
}

func (c *CoapResponse) GetMessage() (*Message) {
    return c.msg
}

func (c *CoapResponse) GetError() (error) {
    return c.err
}