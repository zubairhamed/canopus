package goap

import (
	"testing"
	"bytes"
)

func TestInvalidMessage(t *testing.T) {
    _, err := BytesToMessage(make([]byte, 0))
	if err == nil {
		t.Error("Message should be invalid")
	}

	_, err = BytesToMessage(make([]byte, 4));
	if err == nil {
		t.Error("Message should be invalid")
	}
}

func TestMessageConversion(t *testing.T) {
	msg := NewMessageOfType(TYPE_CONFIRMABLE, 0xf0f0)

	// Byte 1
	msg.Token = []byte("abcd1234")
	msg.AddOption(OPTION_CONTENT_FORMAT, MEDIATYPE_APPLICATION_LINK_FORMAT)

	// Convert Message to BYte
	b := MessageToBytes(msg)

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
