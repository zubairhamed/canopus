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

func NewBasicConfirmableMessage() *Message {
	msg := NewMessageOfType(TYPE_CONFIRMABLE, 0xf0f0)
	msg.Code = GET
	msg.Token = []byte("abcd1234")
	msg.SetStringPayload("xxxxx")

	return msg
}
