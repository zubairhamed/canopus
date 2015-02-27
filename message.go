package goap
import (
    "errors"
    "encoding/binary"
    "fmt"
    "strings"
    "bytes"
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

const (
    MEDIATYPE_TEXT_PLAIN                = 0
    MEDIATYPE_APPLICATION_LINK_FORMAT   = 40
    MEDIATYPE_APPLICATION_XML           = 41
    MEDIATYPE_APPLICATION_OCTET_STREAM  = 42
    MEDIATYPE_APPLICATION_EXI           = 47
    MEDIATYPE_APPLICATION_JSON          = 50
)

const (
	COAPCODE_201_CREATED 					= 65
	COAPCODE_202_DELETED					= 66
	COAPCODE_203_VALID						= 67
	COAPCODE_204_CHANGED					= 68
	COAPCODE_205_CONTENT					= 69
	COAPCODE_400_BAD_REQUEST				= 128
	COAPCODE_401_UNAUTHORIZED				= 129
	COAPCODE_402_BAD_OPTION					= 130
	COAPCODE_403_FORBIDDEN					= 131
	COAPCODE_404_NOT_FOUND					= 132
	COAPCODE_405_METHOD_NOT_ALLOWED			= 133
	COAPCODE_406_NOT_ACCEPTABLE				= 134
	COAPCODE_412_PRECONDITION_FAILED		= 140
	COAPCODE_413_REQUEST_ENTITY_TOO_LARGE	= 141
	COAPCODE_415_UNSUPPORTED_CONTENT_FORMAT	= 143
	COAPCODE_500_INTERNAL_SERVER_ERROR		= 160
	COAPCODE_501_NOT_IMPLEMENTED			= 161
	COAPCODE_502_BAD_GATEWAY				= 162
	COAPCODE_503_SERVICE_UNAVAILABLE		= 163
	COAPCODE_504_GATEWAY_TIMEOUT			= 164
	COAPCODE_505_PROXYING_NOT_SUPPORTED		= 165
)

const PAYLOAD_MARKER = 0xff

func NewMessage() *Message {
	return &Message{}
}

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
func BytesToMessage(data []byte) (*Message, error) {
    msg := NewMessage()

    dataLen := len(data)
    if dataLen < 4 {
        return msg, errors.New("Packet length less than 4 bytes")
    }

    ver := data[DATA_HEADER] >> 6
	if ver != 1 {
		return nil, errors.New("Invalid version")
	}

    msg.MessageType = data[DATA_HEADER] >> 4 & 0x03

	tokenLength := data[DATA_HEADER] & 0x0f

	msg.Code = data[DATA_CODE]

    msg.MessageId = binary.BigEndian.Uint16(data[DATA_MSGID_START:DATA_MSGID_END])

    // Token
    if (tokenLength > 0) {
        msg.Token = make([]byte, tokenLength)
        token := data[DATA_TOKEN_START:DATA_TOKEN_START + tokenLength]
        copy (msg.Token, token)
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
    tmp := data[DATA_TOKEN_START + msg.GetTokenLength():]
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
				msg.Options = append(msg.Options, NewOption(optionId, decodeInt(optionValue)))
				break;

				case OPTION_URI_HOST, OPTION_LOCATION_PATH, OPTION_URI_PATH, OPTION_URI_QUERY,
				 	 OPTION_LOCATION_QUERY, OPTION_PROXY_URI, OPTION_PROXY_SCHEME:
				msg.Options = append(msg.Options, NewOption(optionId, string(optionValue)))
				break;

				default:
                fmt.Println("Ignoring unknown option id " + string(optionId))
                break;
            }
            tmp = tmp[optionLength:]
        }
        lastOptionId = optionId
    }
    msg.Payload = tmp

    err := ValidateMessage(msg)

    return msg, err
}

func MessageToBytes(msg *Message) []byte {
	messageId := []byte{ 0, 0 }
	binary.BigEndian.PutUint16(messageId, msg.MessageId)

	buf := bytes.NewBuffer([]byte{})
	buf.Write([]byte{ (1 << 6) | (msg.MessageType << 4) | 0x0f & msg.GetTokenLength()})
	// buf.Write([]byte{ msg.GetCodeClass() << 5 | 0x0f & msg.GetCodeDetail()})
	buf.Write([]byte{ msg.Code })
	buf.Write([]byte{messageId[0]})
	buf.Write([]byte{messageId[1]})
	buf.Write(msg.Token)


	// sort.Sort(&msg.Options)

	lastOptionId := 0
	for _, opt := range msg.Options {
		b := ValueToBytes(opt.Value)
		optCode := opt.Code
		if len(b) >= 15 {
			buf.Write([]byte{ byte(int(optCode) - lastOptionId) << 4 | 15, byte(len(b) - 15), } )
		} else {
			buf.Write([]byte{ byte(int(optCode) - lastOptionId) << 4 | byte(len(b))} )
		}
	}

	if (len(msg.Payload) > 0) {
		buf.Write([]byte{ PAYLOAD_MARKER  })
	}

	buf.Write(msg.Payload)

	return buf.Bytes()
}

/*
prev := 0
	for _, o := range m.opts {
		b := o.toBytes()
		if len(b) >= 15 {
			buf.Write([]byte{
				byte(int(o.ID)-prev)<<4 | 15,
				byte(len(b) - 15),
			})
		} else {
			buf.Write([]byte{byte(int(o.ID)-prev)<<4 | byte(len(b))})
		}
		if int(o.ID)-prev > 15 {
			return nil, ErrOptionGapTooLarge
		}

		buf.Write(b)
		prev = int(o.ID)
	}
 */



func ValidateMessage(msg *Message) error {
    if msg.MessageType > 3 {
        return errors.New("Unknown message type")
    }

    if msg.GetTokenLength() > 8 {
        return errors.New("Invalid Token Length ( > 8)")
    }

	/*
    codeClass := msg.GetCodeClass()
    if codeClass != 0 && codeClass != 2 && codeClass != 4 && codeClass != 5 {
        return errors.New("Unknown Code class")
    }
    */

    return nil
}

type Message struct {
	Method      uint8
	MessageType uint8
	Code		uint8
	MessageId   uint16
	Payload     []byte
	Token       []byte
	Options     []Option
}

func (c Message) GetCodeString() string {
	codeClass := string(c.Code >> 5)
	codeDetail := string(c.Code & 0x1f)

	return codeClass + "." + codeDetail
}


func (c Message) GetMethod() uint8 {
    return (c.Code & 0x1f)
}

func (c Message) GetTokenLength() uint8 {
	return uint8(len(c.Token))
}

func (c Message) GetOptions(id int) []Option {
    var opts []Option
    for _, val := range c.Options {
        if val.Code == id {
            opts = append(opts, val)
        }
    }
    return opts
}

func (c Message) GetOptionsAsString(id int) []string {
    opts := c.GetOptions(id)

    var str []string
    for _, o := range opts {
        str = append(str, o.Value.(string))
    }
    return str
}

func (c Message) GetPath() string {
    opts := c.GetOptionsAsString(OPTION_URI_PATH)

    return strings.Join(opts, "/")
}

func NewOption(optionNumber int, optionValue interface{}) Option{
    return Option{
        Code: optionNumber,
        Value: optionValue,
    }
}

/* Option */
type Option struct {
    Code     int
    Value   interface{}
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

func encodeInt(v uint32) []byte {
	switch {
	case v == 0:
		return nil

	case v < 256:
		return []byte{ byte(v) }

	case v < 65536:
		rv := []byte{0, 0}
		binary.BigEndian.PutUint16(rv, uint16(v))
		return rv

	case v < 16777216:
		rv := []byte{0, 0, 0, 0}
		binary.BigEndian.PutUint32(rv, uint32(v))
		return rv[1:]

	default:
		rv := []byte{0, 0, 0, 0}
		binary.BigEndian.PutUint32(rv, uint32(v))
		return rv
	}
}

func ValueToBytes(value interface {}) []byte {
	var v uint32

	switch i := value.(type) {
	case string:
		return []byte(i)
	case []byte:
		return i
	case byte:
		v = uint32(i)
	case int:
		v = uint32(i)
	case int32:
		v = uint32(i)
	case uint:
		v = uint32(i)
	case uint32:
		v = i
	default:
		break;
	}

	return encodeInt(v)
}

func PayloadAsString(b []byte) string {
    buff := bytes.NewBuffer(b)

    return buff.String()
}
