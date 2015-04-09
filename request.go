package goap

func NewRequest(messageType uint8, messageCode CoapCode, messageId uint16) *CoapRequest {
    msg := NewMessage(messageType, messageCode, messageId)
    msg.Token = []byte(GenerateToken(8))

    return &CoapRequest{
        msg: msg,
        isConfirmable: true,
    }
}

type CoapRequest struct {
    msg             *Message
    uri             string
    isConfirmable   bool
}

func (c *CoapRequest) GetRequestURI() {

}

func (c *CoapRequest) GetOption() {

}

func (c *CoapRequest) GetOptions() {

}

func (c *CoapRequest) GetMessage() {

}

func (c *CoapRequest) GetMethod() {

}

func (c *CoapRequest) GetError() {

}

func (c *CoapRequest) SetStringPayload(payload string) {
    c.msg.Payload = []byte(payload)
}

func (c *CoapRequest) SetRequestURI(uri string) {
    c.msg.AddOptions(NewPathOptions(uri))
}

func (c *CoapRequest) SetConfirmable(con bool) {
    if con {
        c.msg.MessageType = TYPE_CONFIRMABLE
    } else {
        c.msg.MessageType = TYPE_NONCONFIRMABLE
    }
}

func (c *CoapRequest) SetToken(t string) {
    c.msg.Token = []byte(t)
}

func (c *CoapRequest) IncrementMessageId () {
    c.msg.MessageId = c.msg.MessageId+1
}