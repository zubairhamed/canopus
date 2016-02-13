package canopus

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResponse(t *testing.T) {
	msg := NewEmptyMessage(12345)
	msg.SetStringPayload("hello canopus")
	assert.NotNil(t, NewResponseWithMessage(msg))

	response := NewResponse(msg, ErrUnknownCriticalOption)
	assert.NotNil(t, response)
	assert.Equal(t, uint16(12345), response.GetMessage().MessageID)
	assert.Equal(t, ErrUnknownCriticalOption, response.GetError())

}
