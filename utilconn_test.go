package canopus

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSendMessages(t *testing.T) {
	_, err := SendMessageTo(nil, nil, nil)
	assert.NotNil(t, err)
	assert.Equal(t, ERR_NIL_MESSAGE, err)

	_, err = SendMessageTo(NewEmptyMessage(12345), nil, nil)
	assert.NotNil(t, err)
	assert.Equal(t, ERR_NIL_CONN, err)
}
