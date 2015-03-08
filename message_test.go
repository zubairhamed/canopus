package goap

import (
	"bytes"
	"testing"
)

func TestInvalidMessage(t *testing.T) {
	_, err := BytesToMessage(make([]byte, 0))
	if err == nil {
		t.Error("Message should be invalid")
	}

	_, err = BytesToMessage(make([]byte, 4))
	if err == nil {
		t.Error("Message should be invalid")
	}
}

func TestMessageConversion(t *testing.T) {

	msg := NewBasicConfirmableMessage()
	// Byte 1
	msg.AddOption(OPTION_CONTENT_FORMAT, MEDIATYPE_APPLICATION_LINK_FORMAT)

	// Convert Message to BYte
	b, err := MessageToBytes(msg)

	// Reconvert Back Bytes to Message
	newMsg, err := BytesToMessage(b)
	if err != nil {
		t.Error("An error occured converting bytes to message: ", err)
	}

	if newMsg.MessageType != TYPE_CONFIRMABLE {
		t.Error("Type not the same after reconversion")
	}

	if bytes.NewBuffer(newMsg.Token).String() != "abcd1234" {
		t.Error("Token String not the same after reconversion")
	}

	if newMsg.MessageId != 0xf0f0 {
		t.Error("Message ID not the same after reconversion")
	}

	if len(newMsg.Options) == 0 {
		t.Error("Options should not be 0")
	}
}

func TestMessageBadOptions(t *testing.T) {
	testMsg := NewBasicConfirmableMessage()

	// Unknown Critical Option
	unk := OptionCode(99)
	testMsg.AddOption(unk, 0)
	testMsg.AddOption(OPTION_CONTENT_FORMAT, MEDIATYPE_APPLICATION_LINK_FORMAT)

	_, err := MessageToBytes(testMsg)
	if err == nil {
		if err != ERR_UNKNOWN_CRITICAL_OPTION {
			t.Log("Should throw ERR_UNKNOWN_CRITICAL_OPTION")
			t.Fail()
		}
	}
}


func NewBasicConfirmableMessage() (*Message) {
	msg := NewMessageOfType(TYPE_CONFIRMABLE, 0xf0f0)
	msg.Token = []byte("abcd1234")

	return msg
}
