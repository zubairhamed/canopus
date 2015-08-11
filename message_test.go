package canopus

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMessage(t *testing.T) {
	assert.NotNil(t, NewMessage(TYPE_CONFIRMABLE, GET, 12345))
	assert.NotNil(t, NewEmptyMessage(12345))
}

func TestInvalidMessage(t *testing.T) {
	_, err := BytesToMessage(make([]byte, 0))

	assert.NotNil(t, err, "Message should be invalid")

	_, err = BytesToMessage(make([]byte, 4))
	assert.NotNil(t, err, "Message should be invalid")
}

func TestMessagValidation(t *testing.T) {
	// ValidateMessage()
}

func TestMessageConversion(t *testing.T) {

	msg := NewBasicConfirmableMessage()
	// Byte 1
	msg.AddOption(OPTION_CONTENT_FORMAT, MEDIATYPE_APPLICATION_LINK_FORMAT)

	// Convert Message to BYte
	b, err := MessageToBytes(msg)

	// Reconvert Back Bytes to Message
	newMsg, err := BytesToMessage(b)
	assert.Nil(t, err, "An error occured converting bytes to message")

	assert.Equal(t, 0, int(newMsg.MessageType))	// 0x0: Type Confirmable
	assert.Equal(t, "abcd1234", bytes.NewBuffer(newMsg.Token).String(), "Token String not the same after reconversion")
	assert.Equal(t, 61680, int(newMsg.MessageId), "Message ID not the same after reconversion")
	assert.NotEqual(t, 0, len(newMsg.Options), "Options should not be 0")
}

func TestMessageBadOptions(t *testing.T) {
//	testMsg := NewBasicConfirmableMessage()

	// Unknown Critical Option
//	unk := OptionCode(99)
//	testMsg.AddOption(unk, 0)
//	testMsg.AddOption(OPTION_CONTENT_FORMAT, MEDIATYPE_APPLICATION_LINK_FORMAT)
//	_, err := MessageToBytes(testMsg)
//	assert.NotNil(t, err, "Should throw ERR_UNKNOWN_CRITICAL_OPTION")
}

func TestMessageObject(t *testing.T) {
	msg := &Message{}

	assert.Equal(t, 0, len(msg.Options))
	msg.AddOptions(NewPathOptions("/example"))
	msg.AddOption(OPTION_ACCEPT, MEDIATYPE_APPLICATION_XML)
	msg.AddOption(OPTION_CONTENT_FORMAT, MEDIATYPE_APPLICATION_JSON)
	assert.Equal(t, 3, len(msg.Options))

	opt := msg.GetOption(OPTION_ACCEPT)
	assert.NotNil(t, opt)

	msg.RemoveOptions(OPTION_URI_PATH)
	assert.Equal(t, 2, len(msg.Options))
}

func TestOptionConversion(t *testing.T) {
	preMsg := NewBasicConfirmableMessage()

	// preMsg.Code = TYPE_NONCONFIRMABLE

	preMsg.AddOption(OPTION_IF_MATCH, "")
	preMsg.AddOptions(NewPathOptions("/test"))
	preMsg.AddOption(OPTION_ETAG, "1234567890")
	preMsg.AddOption(OPTION_IF_NONE_MATCH, nil)
	preMsg.AddOption(OPTION_OBSERVE, 0)
	preMsg.AddOption(OPTION_URI_PORT, 1234)
	preMsg.AddOption(OPTION_LOCATION_PATH, "/aaa")
	preMsg.AddOption(OPTION_CONTENT_FORMAT, 1)
	preMsg.AddOption(OPTION_MAX_AGE, 1)
	preMsg.AddOption(OPTION_PROXY_URI, "http://www.google.com")
	preMsg.AddOption(OPTION_PROXY_SCHEME, "http://proxy.scheme")

	converted, _ := MessageToBytes(preMsg)

	postMsg, _ := BytesToMessage(converted)

	PrintMessage(postMsg)
}

func TestNewMessageHelpers(t *testing.T) {
	var msg *Message
	var messageId uint16 = 12345

	msg = EmptyMessage(messageId)
	assert.NotNil(t, msg)
	assert.Equal(t, COAPCODE_0_EMPTY, msg.Code)

	msg = CreatedMessage(messageId)
	assert.NotNil(t, msg)
	assert.Equal(t, COAPCODE_201_CREATED, msg.Code)

	msg = DeletedMessage(messageId)
	assert.NotNil(t, msg)
	assert.Equal(t, COAPCODE_202_DELETED, msg.Code)

	msg = ValidMessage(messageId)
	assert.NotNil(t, msg)
	assert.Equal(t, COAPCODE_203_VALID, msg.Code)

	msg = ChangedMessage(messageId)
	assert.NotNil(t, msg)
	assert.Equal(t, COAPCODE_204_CHANGED, msg.Code)

	msg = ContentMessage(messageId)
	assert.NotNil(t, msg)
	assert.Equal(t, COAPCODE_205_CONTENT, msg.Code)

	msg = BadRequestMessage(messageId)
	assert.NotNil(t, msg)
	assert.Equal(t, COAPCODE_400_BAD_REQUEST, msg.Code)

	msg = UnauthorizedMessage(messageId)
	assert.NotNil(t, msg)
	assert.Equal(t, COAPCODE_401_UNAUTHORIZED, msg.Code)

	msg = BadOptionMessage(messageId)
	assert.NotNil(t, msg)
	assert.Equal(t, COAPCODE_402_BAD_OPTION, msg.Code)

	msg = ForbiddenMessage(messageId)
	assert.NotNil(t, msg)
	assert.Equal(t, COAPCODE_403_FORBIDDEN, msg.Code)

	msg = NotFoundMessage(messageId)
	assert.NotNil(t, msg)
	assert.Equal(t, COAPCODE_404_NOT_FOUND, msg.Code)

	msg = MethodNotAllowedMessage(messageId)
	assert.NotNil(t, msg)
	assert.Equal(t, COAPCODE_405_METHOD_NOT_ALLOWED, msg.Code)

	msg = NotAcceptableMessage(messageId)
	assert.NotNil(t, msg)
	assert.Equal(t, COAPCODE_406_NOT_ACCEPTABLE, msg.Code)

	msg = ConflictMessage(messageId)
	assert.NotNil(t, msg)
	assert.Equal(t, COAPCODE_409_CONFLICT, msg.Code)

	msg = PreconditionFailedMessage(messageId)
	assert.NotNil(t, msg)
	assert.Equal(t, COAPCODE_412_PRECONDITION_FAILED, msg.Code)

	msg = RequestEntityTooLargeMessage(messageId)
	assert.NotNil(t, msg)
	assert.Equal(t, COAPCODE_413_REQUEST_ENTITY_TOO_LARGE, msg.Code)

	msg = UnsupportedContentFormatMessage(messageId)
	assert.NotNil(t, msg)
	assert.Equal(t, COAPCODE_415_UNSUPPORTED_CONTENT_FORMAT, msg.Code)

	msg = InternalServerErrorMessage(messageId)
	assert.NotNil(t, msg)
	assert.Equal(t, COAPCODE_500_INTERNAL_SERVER_ERROR, msg.Code)

	msg = NotImplementedMessage(messageId)
	assert.NotNil(t, msg)
	assert.Equal(t, COAPCODE_501_NOT_IMPLEMENTED, msg.Code)

	msg = BadGatewayMessage(messageId)
	assert.NotNil(t, msg)
	assert.Equal(t, COAPCODE_502_BAD_GATEWAY, msg.Code)

	msg = ServiceUnavailableMessage(messageId)
	assert.NotNil(t, msg)
	assert.Equal(t, COAPCODE_503_SERVICE_UNAVAILABLE, msg.Code)

	msg = GatewayTimeoutMessage(messageId)
	assert.NotNil(t, msg)
	assert.Equal(t, COAPCODE_504_GATEWAY_TIMEOUT, msg.Code)

	msg = ProxyingNotSupportedMessage(messageId)
	assert.NotNil(t, msg)
	assert.Equal(t, COAPCODE_505_PROXYING_NOT_SUPPORTED, msg.Code)
}

func NewBasicConfirmableMessage() *Message {
	msg := NewMessageOfType(TYPE_CONFIRMABLE, 0xf0f0)
	msg.Code = GET
	msg.Token = []byte("abcd1234")
	msg.SetStringPayload("xxxxx")

	return msg
}
