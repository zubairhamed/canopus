package canopus

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
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
	var messageID uint16 = 12345

	testData := []struct {
		msg  *Message
		code CoapCode
	}{
		{EmptyMessage(messageID, MessageAcknowledgment), CoapCodeEmpty},
		{CreatedMessage(messageID, MessageAcknowledgment), CoapCodeCreated},
		{DeletedMessage(messageID, MessageAcknowledgment), CoapCodeDeleted},
		{ValidMessage(messageID, MessageAcknowledgment), CoapCodeValid},
		{ChangedMessage(messageID, MessageAcknowledgment), CoapCodeChanged},
		{ContentMessage(messageID, MessageAcknowledgment), CoapCodeContent},
		{BadRequestMessage(messageID, MessageAcknowledgment), CoapCodeBadRequest},
		{UnauthorizedMessage(messageID, MessageAcknowledgment), CoapCodeUnauthorized},
		{BadOptionMessage(messageID, MessageAcknowledgment), CoapCodeBadOption},
		{ForbiddenMessage(messageID, MessageAcknowledgment), CoapCodeForbidden},
		{NotFoundMessage(messageID, MessageAcknowledgment), CoapCodeNotFound},
		{MethodNotAllowedMessage(messageID, MessageAcknowledgment), CoapCodeMethodNotAllowed},
		{NotAcceptableMessage(messageID, MessageAcknowledgment), CoapCodeNotAcceptable},
		{ConflictMessage(messageID, MessageAcknowledgment), CoapCodeConflict},
		{PreconditionFailedMessage(messageID, MessageAcknowledgment), CoapCodePreconditionFailed},
		{RequestEntityTooLargeMessage(messageID, MessageAcknowledgment), CoapCodeRequestEntityTooLarge},
		{UnsupportedContentFormatMessage(messageID, MessageAcknowledgment), CoapCodeUnsupportedContentFormat},
		{InternalServerErrorMessage(messageID, MessageAcknowledgment), CoapCodeInternalServerError},
		{NotImplementedMessage(messageID, MessageAcknowledgment), CoapCodeNotImplemented},
		{BadGatewayMessage(messageID, MessageAcknowledgment), CoapCodeBadGateway},
		{ServiceUnavailableMessage(messageID, MessageAcknowledgment), CoapCodeServiceUnavailable},
		{GatewayTimeoutMessage(messageID, MessageAcknowledgment), CoapCodeGatewayTimeout},
		{ProxyingNotSupportedMessage(messageID, MessageAcknowledgment), CoapCodeProxyingNotSupported},
	}

	for _, td := range testData {
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
