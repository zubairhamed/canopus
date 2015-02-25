package goap
import (
    "errors"
    "encoding/binary"
    "fmt"
    "strings"
)

const (
    TYPE_CONFIRMABLE        = 0
    TYPE_NONCONFIRMABLE     = 1
    TYPE_ACKNOWLEDGEMENT    = 2
    TYPE_RESET              = 3
)

const (
    METHOD_GET      = 1
    METHOD_POST     = 2
    METHOD_PUT      = 3
    METHOD_DELETE   = 4
)

const (
    DATA_HEADER         = 0
    DATA_CODE           = 1
    DATA_MSGID_START    = 2
    DATA_MSGID_END      = 4
    DATA_TOKEN_START    = 4
)

const (
    OPTION_IF_MATCH         = 1
    OPTION_URI_HOST         = 3
    OPTION_ETAG             = 4
    OPTION_IF_NONE_MATCH    = 5
    OPTION_URI_PORT         = 7
    OPTION_LOCATION_PATH    = 8
    OPTION_URI_PATH         = 11
    OPTION_CONTENT_FORMAT   = 12
    OPTION_MAX_AGE          = 14
    OPTION_URI_QUERY        = 15
    OPTION_ACCEPT           = 17
    OPTION_LOCATION_QUERY   = 20
    OPTION_PROXY_URI        = 35
    OPTION_PROXY_SCHEME     = 39
    OPTION_SIZE1            = 60
)

const PAYLOAD_MARKER = 0xff

/*
     0                   1                   2                   3
    0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |Ver| T |  TKL  |      Code     |          Message ID           |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |   Token (if any, TKL bytes) ...
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |   Options (if any) ...
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |1 1 1 1 1 1 1 1|    Payload (if any) ...
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
*/
func NewMessage(data []byte) (Message, error) {
    msg := &CoApMessage{}

    dataLen := len(data)
    if dataLen < 4 {
        return msg, errors.New("Packet length less than 4 bytes")
    }

    msg.version = data[DATA_HEADER] >> 6
    msg.messageType = data[DATA_HEADER] >> 4 & 0x03
    msg.tokenLength = data[DATA_HEADER] & 0x0f
    msg.codeClass = data[DATA_CODE] >> 5
    msg.codeDetail = data[DATA_CODE] & 0x1f

    msg.messageId = binary.BigEndian.Uint16(data[DATA_MSGID_START:DATA_MSGID_END])

    // Token
    if (msg.TokenLength() > 0) {
        msg.token = make([]byte, msg.TokenLength())
        token := data[DATA_TOKEN_START:DATA_TOKEN_START + msg.TokenLength()]
        copy (msg.token, token)
    }


    /*
    0   1   2   3   4   5   6   7
   +---------------+---------------+
   |               |               |
   |  Option Delta | Option Length |   1 byte
   |               |               |
   +---------------+---------------+
   \                               \
   /         Option Delta          /   0-2 bytes
   \          (extended)           \
   +-------------------------------+
   \                               \
   /         Option Length         /   0-2 bytes
   \          (extended)           \
   +-------------------------------+
   \                               \
   /                               /
   \                               \
   /         Option Value          /   0 or more bytes
   \                               \
   /                               /
   \                               \
   +-------------------------------+
   */
    tmp := data[DATA_TOKEN_START + msg.TokenLength():]
    lastOptionId := 0
    for len(tmp) > 0 {
        if tmp[0] == PAYLOAD_MARKER {
            tmp = tmp[1:]
            break
        }

        optionId := lastOptionId

        optionDelta := int(tmp[0] >> 4)
        optionLength := int(tmp[0] &0x0f)

        tmp = tmp[1:]
        if optionDelta < 13 {
            optionId += int(optionDelta)
        } else {
            switch optionDelta {
                case 13:
                optionDeltaExtended := int(tmp[0])
                optionId += optionDeltaExtended - 13
                tmp = tmp[1:]
                break

                case 14:
                optionDeltaExtended := decodeInt(tmp[:1])
                optionId += int(optionDeltaExtended - uint32(269))
                tmp = tmp[2:]
                break
            }
        }

        if optionLength >= 13 {
            switch optionLength {
                case 13:
                optionLength = int(tmp[0] - 13)
                tmp = tmp[1:]
                break

                case 14:
                optionLength = int(decodeInt(tmp[:1]) - uint32(269))
                tmp = tmp[2:]
                break

                case 15:
                return msg, errors.New("Message Format Error. Option length has reserved value 15")
            }
        }

        if optionLength > 0 {
            optionValue := tmp[:optionLength]

            switch optionId {
                case OPTION_URI_PORT, OPTION_CONTENT_FORMAT, OPTION_MAX_AGE, OPTION_ACCEPT, OPTION_SIZE1:
                msg.AddOption(NewOption(optionId, decodeInt(optionValue)))
                break;

                case OPTION_URI_HOST, OPTION_LOCATION_PATH, OPTION_URI_PATH, OPTION_URI_QUERY,
                     OPTION_LOCATION_QUERY, OPTION_PROXY_URI, OPTION_PROXY_SCHEME:
                msg.AddOption(NewOption(optionId, string(optionValue)))
                break;

                default:
                fmt.Println("Ignoring unknown option id " + string(optionId))
                break;
            }
            tmp = tmp[optionLength:]
        }
        lastOptionId = optionId
    }
    msg.payload = tmp

    err := ValidateMessage(msg)

    return msg, err
}

func ValidateMessage(msg Message) error {
    if msg.Version() != 1 {
        return errors.New("Invalid version")
    }

    if msg.Type() > 3 {
        return errors.New("Unknown message type")
    }

    if msg.TokenLength() > 8 {
        return errors.New("Invalid Token Length ( > 8)")
    }

    codeClass := msg.CodeClass()
    if codeClass != 0 && codeClass != 2 && codeClass != 4 && codeClass != 5 {
        return errors.New("Unknown Code class")
    }

    return nil
}

type Message interface {
    Version() uint8
    Type() uint8
    CodeClass() uint8
    CodeDetail() uint8
    Code() string
    MessageId() uint16
    Method() uint8
    Path() string
    Payload() []byte
    TokenLength() uint8
    Token() []byte
    Options(int) []Option
    OptionsAsString(int) []string
}

type CoApMessage struct {
    method      uint8
    version     uint8
    messageType uint8
    codeClass   uint8
    codeDetail  uint8
    messageId   uint16
    payload     []byte
    tokenLength uint8
    token       []byte
    options     []Option
}

func (c *CoApMessage) Version() uint8 {
    return c.version
}

func (c *CoApMessage) Token() []byte {
    return c.token
}

func (c *CoApMessage) Type() uint8 {
    return c.messageType
}

func (c *CoApMessage) CodeClass() uint8 {
    return c.codeClass
}

func (c *CoApMessage) CodeDetail() uint8 {
    return c.codeDetail
}

func (c *CoApMessage) Code() string {
    return string(c.codeClass) + "." + string(c.codeDetail)
}

func (c *CoApMessage) MessageId() uint16 {
    return c.messageId
}

func (c *CoApMessage) Method() uint8 {
    return c.codeDetail
}

func (c *CoApMessage) Payload() []byte {
    return c.payload
}

func (c *CoApMessage) TokenLength() uint8 {
    return c.tokenLength
}

func (c *CoApMessage) Options(id int) []Option {
    var opts []Option
    for _, val := range c.options {
        if val.num == id {
            opts = append(opts, val)
        }
    }
    return opts
}

func (c *CoApMessage) OptionsAsString(id int) []string {
    opts := c.Options(id)

    var str []string
    for _, o := range opts {
        str = append(str, o.value.(string))
    }
    return str
}

func (c *CoApMessage) Path() string {
    opts := c.OptionsAsString(OPTION_URI_PATH)

    return strings.Join(opts, "/")
}


func (c *CoApMessage) AddOption(o Option) {
    c.options = append(c.options, o)
}

func NewOption(optionNumber int, optionValue interface{}) Option{
    return Option{
        num: optionNumber,
        value: optionValue,
    }
}

/* Option */
type Option struct {
    num     int
    value   interface{}
}

func (o *Option) Name() string {
    return "Name of option"
}


/* Helpers */
func decodeInt(b []byte) uint32 {
    tmp := []byte{0, 0, 0, 0}
    copy(tmp[4-len(b):], b)
    return binary.BigEndian.Uint32(tmp)
}