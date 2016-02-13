package canopus

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMessage(t *testing.T) {
	assert.NotNil(t, NewMessage(MessageConfirmable, Get, 12345))
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
	msg.AddOption(OptionContentFormat, MediaTypeApplicationLinkFormat)

	// Convert Message to BYte
	b, err := MessageToBytes(msg)

	// Reconvert Back Bytes to Message
	newMsg, err := BytesToMessage(b)
	assert.Nil(t, err, "An error occured converting bytes to message")

	assert.Equal(t, 0, int(newMsg.MessageType)) // 0x0: Type Confirmable
	assert.Equal(t, "abcd1234", bytes.NewBuffer(newMsg.Token).String(), "Token String not the same after reconversion")
	assert.Equal(t, 61680, int(newMsg.MessageID), "Message ID not the same after reconversion")
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
	msg.AddOption(OptionAccept, MediaTypeApplicationXML)
	msg.AddOption(OptionContentFormat, MediaTypeApplicationJSON)
	assert.Equal(t, 3, len(msg.Options))

	opt := msg.GetOption(OptionAccept)
	assert.NotNil(t, opt)

	msg.RemoveOptions(OptionURIPath)
	assert.Equal(t, 2, len(msg.Options))
}

func TestOptionConversion(t *testing.T) {
	preMsg := NewBasicConfirmableMessage()

	preMsg.AddOption(OptionIfMatch, "")
	preMsg.AddOptions(NewPathOptions("/test"))
	preMsg.AddOption(OptionEtag, "1234567890")
	preMsg.AddOption(OptionIfNoneMatch, nil)
	preMsg.AddOption(OptionObserve, 0)
	preMsg.AddOption(OptionURIPort, 1234)
	preMsg.AddOption(OptionLocationPath, "/aaa")
	preMsg.AddOption(OptionContentFormat, 1)
	preMsg.AddOption(OptionMaxAge, 1)
	preMsg.AddOption(OptionProxyURI, "http://www.google.com")
	preMsg.AddOption(OptionProxyScheme, "http://proxy.scheme")

	converted, _ := MessageToBytes(preMsg)

	postMsg, _ := BytesToMessage(converted)

	PrintMessage(postMsg)
}

func TestNewMessageHelpers(t *testing.T) {
	var messageId uint16 = 12345

	test_data := []struct {
		msg  *Message
		code CoapCode
	}{
		{EmptyMessage(messageId, MessageAcknowledgment), CoapCodeEmpty},
		{CreatedMessage(messageId, MessageAcknowledgment), CoapCodeCreated},
		{DeletedMessage(messageId, MessageAcknowledgment), CoapCodeDeleted},
		{ValidMessage(messageId, MessageAcknowledgment), CoapCodeValid},
		{ChangedMessage(messageId, MessageAcknowledgment), CoapCodeChanged},
		{ContentMessage(messageId, MessageAcknowledgment), CoapCodeContent},
		{BadRequestMessage(messageId, MessageAcknowledgment), CoapCodeBadRequest},
		{UnauthorizedMessage(messageId, MessageAcknowledgment), CoapCodeUnauthorized},
		{BadOptionMessage(messageId, MessageAcknowledgment), CoapCodeBadOption},
		{ForbiddenMessage(messageId, MessageAcknowledgment), CoapCodeForbidden},
		{NotFoundMessage(messageId, MessageAcknowledgment), CoapCodeNotFound},
		{MethodNotAllowedMessage(messageId, MessageAcknowledgment), CoapCodeMethodNotAllowed},
		{NotAcceptableMessage(messageId, MessageAcknowledgment), CoapCodeNotAcceptable},
		{ConflictMessage(messageId, MessageAcknowledgment), CoapCodeConflict},
		{PreconditionFailedMessage(messageId, MessageAcknowledgment), CoapCodePreconditionFailed},
		{RequestEntityTooLargeMessage(messageId, MessageAcknowledgment), CoapCodeRequestEntityTooLarge},
		{UnsupportedContentFormatMessage(messageId, MessageAcknowledgment), CoapCodeUnsupportedContentFormat},
		{InternalServerErrorMessage(messageId, MessageAcknowledgment), CoapCodeInternalServerError},
		{NotImplementedMessage(messageId, MessageAcknowledgment), CoapCodeNotImplemented},
		{BadGatewayMessage(messageId, MessageAcknowledgment), CoapCodeBadGateway},
		{ServiceUnavailableMessage(messageId, MessageAcknowledgment), CoapCodeServiceUnavailable},
		{GatewayTimeoutMessage(messageId, MessageAcknowledgment), CoapCodeGatewayTimeout},
		{ProxyingNotSupportedMessage(messageId, MessageAcknowledgment), CoapCodeProxyingNotSupported},
	}

	for _, td := range test_data {
		assert.NotNil(t, td.msg)
		assert.Equal(t, td.code, td.msg.Code)
	}
}

func NewBasicConfirmableMessage() *Message {
	msg := NewMessageOfType(MessageConfirmable, 0xf0f0)
	msg.Code = Get
	msg.Token = []byte("abcd1234")
	msg.SetStringPayload("xxxxx")

	return msg
}
