package canopus
import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestRoutes(t *testing.T) {
	re, static := CreateCompilableRoutePath("/")
	assert.True(t, static)
	assert.True(t, re.MatchString("/"))

	re, static = CreateCompilableRoutePath("/test")
	assert.True(t, static)
	assert.False(t, re.MatchString("/"))
	assert.True(t, re.MatchString("/test"))

	re, static = CreateCompilableRoutePath("/test/:var")
	assert.False(t, static)
	assert.True(t, re.MatchString("/test/abc"))
	assert.False(t, re.MatchString("/test/abc/def"))

	re, static = CreateCompilableRoutePath("/test.abc/:var")
	assert.False(t, static)
	assert.True(t, re.MatchString("/test.abc/abc"))
	assert.False(t, re.MatchString("/test.abc/abc/def"))

}