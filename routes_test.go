package canopus

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoutes(t *testing.T) {
	var route Route
	var matches bool

	route = CreateNewRegExRoute("/", "GET", nil)
	matches, _ = route.Matches("/")
	assert.True(t, matches)

	route = CreateNewRegExRoute("/test", "GET", nil)
	matches, _ = route.Matches("/")
	assert.False(t, matches)
	matches, _ = route.Matches("/test")
	assert.True(t, matches)

	route = CreateNewRegExRoute("/test/:var", "GET", nil)
	matches, _ = route.Matches("/test/abc")
	assert.True(t, matches)
	matches, _ = route.Matches("/test/abc/def")
	assert.False(t, matches)

	route = CreateNewRegExRoute("/test/:var/foo", "GET", nil)
	matches, _ = route.Matches("/test/abc/foo")
	assert.True(t, matches)
	matches, _ = route.Matches("/test/abc")
	assert.False(t, matches)
	matches, _ = route.Matches("/test/abc/def")
	assert.False(t, matches)
	matches, _ = route.Matches("/test//foo")
	assert.False(t, matches)
	matches, _ = route.Matches("/test/foo")
	assert.False(t, matches)

	route = CreateNewRegExRoute("/test.abc/:var", "GET", nil)
	matches, _ = route.Matches("/test.abc/abc")
	assert.True(t, matches)
	matches, _ = route.Matches("/test.abc/abc/def")
	assert.False(t, matches)
}
