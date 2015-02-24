package goap

type MessageType    uint8
type Method         uint8
type MessageCode    uint8
type MessageId      uint16
type MessagePayload []byte

const (
    TYPE_CONFIRMABLE        MessageType = 0
    TYPE_NONCONFIRMABLE     MessageType = 1
    TYPE_ACKNOWLEDGEMENT    MessageType = 2
    TYPE_RESET              MessageType = 3
)

const (
    METHOD_GET      Method = 1
    METHOD_POST     = 2
    METHOD_PUT      = 3
    METHOD_DELETE   = 4
)

// Message
func NewMessage() Message {
    return &CoApMessage{}
}

type Message interface {
    Version() uint8
    Type() MessageType
    Code() MessageCode
    MessageId() MessageId
    Path() string
    Method() Method
    Payload() MessagePayload
}

type CoApMessage struct {
    path        string
    method      Method
    version     uint8
    messageType MessageType
    code        MessageCode
    messageId   MessageId
    payload     MessagePayload
}

func (c *CoApMessage) Version() uint8 {
    return c.version
}

func (c *CoApMessage) Type() MessageType {
    return c.messageType
}

func (c *CoApMessage) Code() MessageCode {
    return c.code
}

func (c *CoApMessage) MessageId() MessageId {
    return c.messageId
}

func (c *CoApMessage) Path() string {
    return c.path
}

func (c *CoApMessage) Method() Method {
    return c.method
}

func (c *CoApMessage) Payload() MessagePayload {
    return c.payload
}

// Functions
func IsRequestMessage(m Message) bool {
    if m.Code() >= 1 && m.Code() <= 31 {
        return true
    }
    return false
}

func IsResponseMessage(m Message) bool {
    if m.Code() >= 64 && m.Code() <= 191 {
        return true
    }
    return false
}

func IsEmptyMessage(m Message) bool {
    if m.Code() == 0 {
        return true
    }
    return false
}

func ParseMessage(data []byte) Message {
    return NewMessage()
}