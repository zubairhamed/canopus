package canopus

import (
	"bytes"
	"encoding/binary"
	"errors"
	"log"
	"sort"
	"strconv"
	"strings"
)

// Instantiates a new message object
// messageType (e.g. Confirm/Non-Confirm)
// CoAP code	404 - Not found etc
// Message ID	uint16 unique id
func NewMessage(messageType uint8, code CoapCode, messageID uint16) *Message {
	return &Message{
		MessageType: messageType,
		MessageID:   messageID,
		Code:        code,
	}
}

// Instantiates an empty message with a given message id
func NewEmptyMessage(id uint16) *Message {
	msg := NewMessageOfType(MessageAcknowledgment, id)

	return msg
}

// Instantiates an empty message of a specific type and message id
func NewMessageOfType(t uint8, id uint16) *Message {
	return &Message{
		MessageType: t,
		MessageID:   id,
	}
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

// Converts an array of bytes to a Mesasge object.
// An error is returned if a parsing error occurs
func BytesToMessage(data []byte) (*Message, error) {
	msg := &Message{}

	dataLen := len(data)
	if dataLen < 4 {
		return msg, ErrPacketLengthLessThan4
	}

	ver := data[DataHeader] >> 6
	if ver != 1 {
		return nil, ErrInvalidCoapVersion
	}

	msg.MessageType = data[DataHeader] >> 4 & 0x03
	tokenLength := data[DataHeader] & 0x0f
	msg.Code = CoapCode(data[DataCode])

	msg.MessageID = binary.BigEndian.Uint16(data[DataMsgIDStart:DataMsgIDEnd])

	// Token
	if tokenLength > 0 {
		msg.Token = make([]byte, tokenLength)
		token := data[DataTokenStart : DataTokenStart+tokenLength]
		copy(msg.Token, token)
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
	tmp := data[DataTokenStart+msg.GetTokenLength():]

	lastOptionID := 0
	for len(tmp) > 0 {
		if tmp[0] == PayloadMarker {
			tmp = tmp[1:]
			break
		}

		optionDelta := int(tmp[0] >> 4)
		optionLength := int(tmp[0] & 0x0f)

		tmp = tmp[1:]
		switch optionDelta {
		case 13:
			optionDeltaExtended := int(tmp[0])
			optionDelta += optionDeltaExtended
			tmp = tmp[1:]
			break

		case 14:
			optionDeltaExtended := decodeInt(tmp[:1])
			optionDelta += int(optionDeltaExtended - uint32(269))
			tmp = tmp[2:]
			break

		case 15:
			return msg, ErrOptionDeltaUsesValue15
		}
		lastOptionID += optionDelta

		switch optionLength {
		case 13:
			optionLengthExtended := int(tmp[0])
			optionLength += optionLengthExtended
			tmp = tmp[1:]
			break

		case 14:
			optionLengthExtended := decodeInt(tmp[:1])
			optionLength += int(optionLengthExtended - uint32(269))
			tmp = tmp[2:]
			break

		case 15:
			return msg, ErrOptionLengthUsesValue15
		}

		optCode := OptionCode(lastOptionID)
		if optionLength > 0 {
			optionValue := tmp[:optionLength]

			switch optCode {
			case OptionURIPort, OptionContentFormat, OptionMaxAge, OptionAccept, OptionSize1,
				OptionSize2, OptionBlock1, OptionBlock2:
				msg.Options = append(msg.Options, NewOption(optCode, decodeInt(optionValue)))
				break

			case OptionURIHost, OptionEtag, OptionLocationPath, OptionURIPath, OptionURIQuery,
				OptionLocationQuery, OptionProxyURI, OptionProxyScheme, OptionObserve:
				msg.Options = append(msg.Options, NewOption(optCode, string(optionValue)))
				break

			default:
				if lastOptionID&0x01 == 1 {
					log.Println("Unknown Critical Option id " + strconv.Itoa(lastOptionID))
					return msg, ErrUnknownCriticalOption
				}
				log.Println("Unknown Option id " + strconv.Itoa(lastOptionID))
				break
			}
			tmp = tmp[optionLength:]
		} else {
			msg.Options = append(msg.Options, NewOption(optCode, nil))
		}
	}
	msg.Payload = NewBytesPayload(tmp)
	err := ValidateMessage(msg)

	return msg, err
}

// type to sort the coap options list (which is mandatory) prior to transmission
type SortOptions []Option

func (opts SortOptions) Len() int {
	return len(opts)
}

func (opts SortOptions) Swap(i, j int) {
	opts[i], opts[j] = opts[j], opts[i]
}

func (opts SortOptions) Less(i, j int) bool {
	return opts[i].GetCode() < opts[j].GetCode()
}

// Converts a message object to a byte array. Typically done prior to transmission
func MessageToBytes(msg *Message) ([]byte, error) {
	messageID := []byte{0, 0}
	binary.BigEndian.PutUint16(messageID, msg.MessageID)

	buf := bytes.Buffer{}
	buf.Write([]byte{(1 << 6) | (msg.MessageType << 4) | 0x0f&msg.GetTokenLength()})
	buf.Write([]byte{byte(msg.Code)})
	buf.Write([]byte{messageID[0]})
	buf.Write([]byte{messageID[1]})
	buf.Write(msg.Token)

	// Sort Options
	sort.Sort(SortOptions(msg.Options))

	lastOptionCode := 0
	for _, opt := range msg.Options {
		optCode := int(opt.GetCode())
		optDelta := optCode - lastOptionCode
		optDeltaValue, _ := getOptionHeaderValue(optDelta)

		byteValue := valueToBytes(opt.GetValue())
		valueLength := len(byteValue)
		optLength := valueLength
		optLengthValue, _ := getOptionHeaderValue(optLength)

		buf.Write([]byte{byte(optDeltaValue<<4 | optLengthValue)})

		if optDeltaValue == 13 {
			buf.Write([]byte{byte(optDelta - 13)})
		} else if optDeltaValue == 14 {
			tmpBuf := new(bytes.Buffer)
			binary.Write(tmpBuf, binary.BigEndian, uint16(optDelta-269))
			buf.Write(tmpBuf.Bytes())
		}

		if optLengthValue == 13 {
			buf.Write([]byte{byte(optLength - 13)})
		} else if optLengthValue == 14 {
			tmpBuf := new(bytes.Buffer)
			binary.Write(tmpBuf, binary.BigEndian, uint16(optLength-269))
			buf.Write(tmpBuf.Bytes())
		}

		buf.Write(byteValue)
		lastOptionCode = int(optCode)
	}

	if msg.Payload != nil {
		if msg.Payload.Length() > 0 {
			buf.Write([]byte{PayloadMarker})
		}
		buf.Write(msg.Payload.GetBytes())
	}
	return buf.Bytes(), nil
}

func getOptionHeaderValue(optValue int) (int, error) {
	switch true {
	case optValue <= 12:
		return optValue, nil

	case optValue <= 268:
		return 13, nil

	case optValue <= 65804:
		return 14, nil
	}
	return 0, errors.New("Invalid Option Delta")
}

// Validates a message object and returns any error upon validation failure
func ValidateMessage(msg *Message) error {
	if msg.MessageType > 3 {
		return ErrUnknownMessageType
	}

	if msg.GetTokenLength() > 8 {
		return ErrInvalidTokenLength
	}

	// Repeated Unrecognized Options
	for _, opt := range msg.Options {
		opts := msg.GetOptions(opt.GetCode())

		if len(opts) > 1 {
			if !IsRepeatableOption(opts[0]) {
				if opts[0].GetCode() & 0x01 == 1 {
					return ErrUnknownCriticalOption
				}
			}
		}
	}

	return nil
}

func NewBlockMessage() *BlockMessage {
	return &BlockMessage{
		Sequence: 0,
	}
}

type BlockMessage struct {
	MessageBuf []byte
	Sequence   uint
}

type BySequence []*BlockMessage

func (o BySequence) Len() int {
	return len(o)
}

func (o BySequence) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}

func (o BySequence) Less(i, j int) bool {
	return o[i].Sequence < o[j].Sequence
}

// A Message object represents a CoAP payload
type Message struct {
	MessageType uint8
	Code        CoapCode
	MessageID   uint16
	Payload     MessagePayload
	Token       []byte
	Options     []Option
}

func (m *Message) GetAcceptedContent() MediaType {
	mediaTypeCode := m.GetOption(OptionAccept).IntValue()

	return MediaType(mediaTypeCode)
}

func (m *Message) GetCodeString() string {
	codeClass := string(m.Code >> 5)
	codeDetail := string(m.Code & 0x1f)

	return codeClass + "." + codeDetail
}

func (m *Message) GetMethod() uint8 {
	return (byte(m.Code) & 0x1f)
}

func (m *Message) GetTokenLength() uint8 {
	return uint8(len(m.Token))
}

func (m *Message) GetTokenString() string {
	return string(m.Token[:len(m.Token)])
}

// Returns an array of options given an option code
func (m Message) GetOptions(id OptionCode) []Option {
	var opts []Option
	for _, val := range m.Options {
		if val.GetCode() == id {
			opts = append(opts, val)
		}
	}
	return opts
}

// Returns the first option found for a given option code
func (m Message) GetOption(id OptionCode) Option {
	for _, val := range m.Options {
		if val.GetCode() == id {
			return val
		}
	}
	return nil
}

// Attempts to return the string value of an Option
func (m Message) GetOptionsAsString(id OptionCode) []string {
	opts := m.GetOptions(id)

	var str []string
	for _, o := range opts {
		if o.GetValue() != nil {
			str = append(str, o.GetValue().(string))
		}
	}
	return str
}

// Returns the string value of the Location Path Options by joining and defining a / separator
func (m *Message) GetLocationPath() string {
	opts := m.GetOptionsAsString(OptionLocationPath)

	return strings.Join(opts, "/")
}

// Returns the string value of the Uri Path Options by joining and defining a / separator
func (m Message) GetURIPath() string {
	opts := m.GetOptionsAsString(OptionURIPath)

	return "/" + strings.Join(opts, "/")
}

// Add an Option to the message. If an option is not repeatable, it will replace
// any existing defined Option of the same type
func (m *Message) AddOption(code OptionCode, value interface{}) {
	opt := NewOption(code, value)
	if IsRepeatableOption(opt) {
		m.Options = append(m.Options, opt)
	} else {
		m.RemoveOptions(code)
		m.Options = append(m.Options, opt)
	}
}

// Add an array of Options to the message. If an option is not repeatable, it will replace
// any existing defined Option of the same type
func (m *Message) AddOptions(opts []Option) {
	for _, opt := range opts {
		if IsRepeatableOption(opt) {
			m.Options = append(m.Options, opt)
		} else {
			m.RemoveOptions(opt.GetCode())
			m.Options = append(m.Options, opt)
		}
	}
}

func (c *Message) SetBlock1Option(opt *Block1Option) {
	c.AddOption(OptionBlock1, opt.GetValue())
}

// Copies the given list of options from another message to this one
func (m *Message) CloneOptions(cm *Message, opts ...OptionCode) {
	for _, opt := range opts {
		m.AddOptions(cm.GetOptions(opt))
	}
}

// Replace an Option
func (m *Message) ReplaceOptions(code OptionCode, opts []Option) {
	m.RemoveOptions(code)

	m.AddOptions(opts)
}

// Removes an Option
func (m *Message) RemoveOptions(id OptionCode) {
	var opts []Option
	for _, opt := range m.Options {
		if opt.GetCode() != id {
			opts = append(opts, opt)
		}
	}
	m.Options = opts
}

// Adds a string payload
func (m *Message) SetStringPayload(s string) {
	m.Payload = NewPlainTextPayload(s)
}

// Determines if a message contains options for proxying (i.e. Proxy-Scheme or Proxy-Uri)
func IsProxyRequest(msg *Message) bool {
	if msg.GetOption(OptionProxyScheme) != nil || msg.GetOption(OptionProxyURI) != nil {
		return true
	}
	return false
}

func valueToBytes(value interface{}) []byte {
	var v uint32

	switch i := value.(type) {
	case string:
		return []byte(i)
	case []byte:
		return i
	case MediaType:
		v = uint32(i)
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
		break
	}

	return encodeInt(v)
}

// Returns the string value for a Message Payload
func PayloadAsString(p MessagePayload) string {
	if p == nil {
		return ""
	}
	return p.String()
}

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
		return []byte{byte(v)}

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

// Determines if a message contains URI targeting a CoAP resource
func IsCoapURI(uri string) bool {
	if strings.HasPrefix(uri, "coap") || strings.HasPrefix(uri, "coaps") {
		return true
	}
	return false
}

// Determines if a message contains URI targeting an HTTP resource
func IsHTTPURI(uri string) bool {
	if strings.HasPrefix(uri, "http") || strings.HasPrefix(uri, "https") {
		return true
	}
	return false

}

// Gets the string representation of a CoAP Method code (e.g. GET, PUT, DELETE etc)
func MethodString(c CoapCode) string {
	switch c {
	case Get:
		return "GET"

	case Delete:
		return "DELETE"

	case Post:
		return "POST"

	case Put:
		return "PUT"
	}
	return ""
}

// Response Code Messages
// Creates a Non-Confirmable Empty Message
func EmptyMessage(messageID uint16, messageType uint8) *Message {
	return NewMessage(messageType, CoapCodeEmpty, messageID)
}

// Creates a Non-Confirmable with CoAP Code 201 - Created
func CreatedMessage(messageID uint16, messageType uint8) *Message {
	return NewMessage(messageType, CoapCodeCreated, messageID)
}

// // Creates a Non-Confirmable with CoAP Code 202 - Deleted
func DeletedMessage(messageID uint16, messageType uint8) *Message {
	return NewMessage(messageType, CoapCodeDeleted, messageID)
}

// Creates a Non-Confirmable with CoAP Code 203 - Valid
func ValidMessage(messageID uint16, messageType uint8) *Message {
	return NewMessage(messageType, CoapCodeValid, messageID)
}

// Creates a Non-Confirmable with CoAP Code 204 - Changed
func ChangedMessage(messageID uint16, messageType uint8) *Message {
	return NewMessage(messageType, CoapCodeChanged, messageID)
}

// Creates a Non-Confirmable with CoAP Code 205 - Content
func ContentMessage(messageID uint16, messageType uint8) *Message {
	return NewMessage(messageType, CoapCodeContent, messageID)
}

// Creates a Non-Confirmable with CoAP Code 400 - Bad Request
func BadRequestMessage(messageID uint16, messageType uint8) *Message {
	return NewMessage(messageType, CoapCodeBadRequest, messageID)
}

func ContinueMessage(messageID uint16, messageType uint8) *Message {
	return NewMessage(messageType, CoapCodeContinue, messageID)
}

// Creates a Non-Confirmable with CoAP Code 401 - Unauthorized
func UnauthorizedMessage(messageID uint16, messageType uint8) *Message {
	return NewMessage(messageType, CoapCodeUnauthorized, messageID)
}

// Creates a Non-Confirmable with CoAP Code 402 - Bad Option
func BadOptionMessage(messageID uint16, messageType uint8) *Message {
	return NewMessage(messageType, CoapCodeBadOption, messageID)
}

// Creates a Non-Confirmable with CoAP Code 403 - Forbidden
func ForbiddenMessage(messageID uint16, messageType uint8) *Message {
	return NewMessage(messageType, CoapCodeForbidden, messageID)
}

// Creates a Non-Confirmable with CoAP Code 404 - Not Found
func NotFoundMessage(messageID uint16, messageType uint8) *Message {
	return NewMessage(messageType, CoapCodeNotFound, messageID)
}

// Creates a Non-Confirmable with CoAP Code 405 - Method Not Allowed
func MethodNotAllowedMessage(messageID uint16, messageType uint8) *Message {
	return NewMessage(messageType, CoapCodeMethodNotAllowed, messageID)
}

// Creates a Non-Confirmable with CoAP Code 406 - Not Acceptable
func NotAcceptableMessage(messageID uint16, messageType uint8) *Message {
	return NewMessage(messageType, CoapCodeNotAcceptable, messageID)
}

// Creates a Non-Confirmable with CoAP Code 409 - Conflict
func ConflictMessage(messageID uint16, messageType uint8) *Message {
	return NewMessage(messageType, CoapCodeConflict, messageID)
}

// Creates a Non-Confirmable with CoAP Code 412 - Precondition Failed
func PreconditionFailedMessage(messageID uint16, messageType uint8) *Message {
	return NewMessage(messageType, CoapCodePreconditionFailed, messageID)
}

// Creates a Non-Confirmable with CoAP Code 413 - Request Entity Too Large
func RequestEntityTooLargeMessage(messageID uint16, messageType uint8) *Message {
	return NewMessage(messageType, CoapCodeRequestEntityTooLarge, messageID)
}

// Creates a Non-Confirmable with CoAP Code 415 - Unsupported Content Format
func UnsupportedContentFormatMessage(messageID uint16, messageType uint8) *Message {
	return NewMessage(messageType, CoapCodeUnsupportedContentFormat, messageID)
}

// Creates a Non-Confirmable with CoAP Code 500 - Internal Server Error
func InternalServerErrorMessage(messageID uint16, messageType uint8) *Message {
	return NewMessage(messageType, CoapCodeInternalServerError, messageID)
}

// Creates a Non-Confirmable with CoAP Code 501 - Not Implemented
func NotImplementedMessage(messageID uint16, messageType uint8) *Message {
	return NewMessage(messageType, CoapCodeNotImplemented, messageID)
}

// Creates a Non-Confirmable with CoAP Code 502 - Bad Gateway
func BadGatewayMessage(messageID uint16, messageType uint8) *Message {
	return NewMessage(messageType, CoapCodeBadGateway, messageID)
}

// Creates a Non-Confirmable with CoAP Code 503 - Service Unavailable
func ServiceUnavailableMessage(messageID uint16, messageType uint8) *Message {
	return NewMessage(messageType, CoapCodeServiceUnavailable, messageID)
}

// Creates a Non-Confirmable with CoAP Code 504 - Gateway Timeout
func GatewayTimeoutMessage(messageID uint16, messageType uint8) *Message {
	return NewMessage(messageType, CoapCodeGatewayTimeout, messageID)
}

// Creates a Non-Confirmable with CoAP Code 505 - Proxying Not Supported
func ProxyingNotSupportedMessage(messageID uint16, messageType uint8) *Message {
	return NewMessage(messageType, CoapCodeProxyingNotSupported, messageID)
}
