package goap

import (
	"testing"
	"bytes"
    "github.com/zubairhamed/goap"
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
	msg := goap.NewMessage()

	// Byte 1
    msg.MessageType = TYPE_CONFIRMABLE
	msg.Token = []byte("abcd1234")

	// Byte 2
    msg.MessageId = 0xf0f0

	b := MessageToBytes(msg)
	newMsg, err := BytesToMessage(b)
	if err != nil {
		t.Error("An error occured converting bytes to message: ", err)
	}

	if newMsg.Type() != TYPE_CONFIRMABLE {
		t.Error("Type not the same after reconversion")
	}

	if bytes.NewBuffer(newMsg.Token()).String() != "abcd1234" {
		t.Error("Token String not the same after reconversion")
	}

	if newMsg.MessageId() != 0xf0f0 {
		t.Error("Message ID not the same after reconversion")
	}
}
