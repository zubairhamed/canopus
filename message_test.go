package goap

import (
	"testing"
	"bytes"
)

func TestInvalidMessage(t *testing.T) {
	_, err := BytesToMessage(make([]byte, 0));
	if err == nil {
		t.Error("Message should be invalid")
	}

	_, err = BytesToMessage(make([]byte, 4));
	if err == nil {
		t.Error("Message should be invalid")
	}
}

func TestMessageConversion(t *testing.T) {
	msg := CoApMessage{}

	// Byte 1
	msg.version = 1
	msg.messageType = TYPE_CONFIRMABLE
	msg.token = []byte("abcd1234")

	// Byte 2
	msg.codeClass = CODECLASS_REQUEST
	msg.codeDetail = METHOD_GET
	msg.messageId = 0xf0f0

	b := MessageToBytes(msg)
	newMsg, err := BytesToMessage(b)
	if err != nil {
		t.Error("An error occured converting bytes to message: ", err)
	}

	if newMsg.Version() != 1 {
		t.Error("Version not equal to 1")
	}

	if newMsg.Type() != TYPE_CONFIRMABLE {
		t.Error("Type not the same after reconversion")
	}

	if bytes.NewBuffer(newMsg.Token()).String() != "abcd1234" {
		t.Error("Token String not the same after reconversion")
	}

	if newMsg.CodeClass() != CODECLASS_REQUEST {
		t.Error("Code Class not the same after reconversion")
	}

	if newMsg.CodeDetail() != METHOD_GET {
		t.Error("Code Detail not the same after reconversion")
	}

	if newMsg.MessageId() != 0xf0f0 {
		t.Error("Message ID not the same after reconversion")
	}
}
