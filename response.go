package goap
import "strings"

func NewResponse(msg *Message, err error) *CoapResponse {
    resp := &CoapResponse{
        msg: msg,
        err: err,
    }

    return resp
}

func NewResponseWithMessage(msg *Message) *CoapResponse {
    resp := &CoapResponse{
        msg: msg,
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

func (c *CoapResponse) GetUriQuery(q string) string {
    qs := c.GetMessage().GetOptionsAsString(OPTION_URI_QUERY)

    for _, o := range qs {
        ps := strings.Split(o, "=")
        if len(ps)  == 2 {
            if ps[0] == q {
                return ps[1]
            }
        }
    }
    return ""
}
