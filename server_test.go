package canopus
import (
	"testing"
	"github.com/stretchr/testify/assert"
)


func TestServerInstantiate(t *testing.T) {
	var s *CoapServer
	s = NewCoapServer("localhost:1000")
	assert.Equal(t, 1000, s.localAddr.Port)
	assert.Equal(t, "udp", s.localAddr.Network())

	s = NewLocalServer()
	assert.NotNil(t, s)
	assert.Equal(t, 5683, s.localAddr.Port)
	assert.Equal(t, "udp", s.localAddr.Network())
}