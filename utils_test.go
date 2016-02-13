package canopus

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateMessageId(t *testing.T) {

	var id, id2 uint16
	id = GenerateMessageID()
	for i := 0; i < 100; i++ {
		id2 = id + 1
		id = GenerateMessageID()
		assert.NotEqual(t, 65535, id)
		assert.Equal(t, id2, id)
	}

	CurrentMessageID = 65535
	id = GenerateMessageID()
	assert.Equal(t, uint16(1), id)
}

func TestGenerateToken(t *testing.T) {
	assert.Equal(t, "", GenerateToken(0))

	for i := 1; i < 10; i++ {
		tok := GenerateToken(i)
		assert.NotEqual(t, "", tok)
		assert.Equal(t, i, len(tok))
	}
}

func TestCoreResourceUtil(t *testing.T) {
	var resources []*CoreResource

	resources = CoreResourcesFromString("")

	assert.Equal(t, 0, len(resources))

	resources = CoreResourcesFromString("</sensors/temp>;ct=41;rt=\"temperature-c\";if=\"sensor\", </sensors/light>;ct=41;rt=\"light-lux\";if=\"sensor\"")

	assert.Equal(t, 2, len(resources))
	resource1 := resources[0]

	assert.Equal(t, "/sensors/temp", resource1.Target)
	assert.Equal(t, 3, len(resource1.Attributes))

	assert.Nil(t, resource1.GetAttribute("invalid_attr"))

	assert.NotNil(t, resource1.GetAttribute("ct"))
	assert.Equal(t, "ct", resource1.GetAttribute("ct").Key)
	assert.Equal(t, "41", resource1.GetAttribute("ct").Value)

	assert.NotNil(t, resource1.GetAttribute("rt"))
	assert.Equal(t, "rt", resource1.GetAttribute("rt").Key)
	assert.Equal(t, "temperature-c", resource1.GetAttribute("rt").Value)

	assert.NotNil(t, resource1.GetAttribute("if"))
	assert.Equal(t, "if", resource1.GetAttribute("if").Key)
	assert.Equal(t, "sensor", resource1.GetAttribute("if").Value)

	resource2 := resources[1]
	assert.Equal(t, "/sensors/light", resource2.Target)
	assert.Equal(t, 3, len(resource2.Attributes))

	assert.NotNil(t, resource2.GetAttribute("ct"))
	assert.Equal(t, "ct", resource2.GetAttribute("ct").Key)
	assert.Equal(t, "41", resource2.GetAttribute("ct").Value)

	assert.NotNil(t, resource2.GetAttribute("rt"))
	assert.Equal(t, "rt", resource2.GetAttribute("rt").Key)
	assert.Equal(t, "light-lux", resource2.GetAttribute("rt").Value)

	assert.NotNil(t, resource2.GetAttribute("if"))
	assert.Equal(t, "if", resource2.GetAttribute("if").Key)
	assert.Equal(t, "sensor", resource2.GetAttribute("if").Value)
}

func TestCoapCodeToString(t *testing.T) {
	testData := []struct {
		coapCode   CoapCode
		codeString string
	}{
		{Get, "GET"},
		{Post, "POST"},
		{Put, "PUT"},
		{Delete, "DELETE"},
		{CoapCodeEmpty, "0 Empty"},
		{CoapCodeCreated, "201 Created"},
		{CoapCodeDeleted, "202 Deleted"},
		{CoapCodeValid, "203 Valid"},
		{CoapCodeChanged, "204 Changed"},
		{CoapCodeContent, "205 Content"},
		{CoapCodeBadRequest, "400 Bad Request"},
		{CoapCodeUnauthorized, "401 Unauthorized"},
		{CoapCodeBadOption, "402 Bad Option"},
		{CoapCodeForbidden, "403 Forbidden"},
		{CoapCodeNotFound, "404 Not Found"},
		{CoapCodeMethodNotAllowed, "405 Method Not Allowed"},
		{CoapCodeNotAcceptable, "406 Not Acceptable"},
		{CoapCodePreconditionFailed, "412 Precondition Failed"},
		{CoapCodeRequestEntityTooLarge, "413 Request Entity Too Large"},
		{CoapCodeUnsupportedContentFormat, "415 Unsupported Content Format"},
		{CoapCodeInternalServerError, "500 Internal Server Error"},
		{CoapCodeNotImplemented, "501 Not Implemented"},
		{CoapCodeBadGateway, "502 Bad Gateway"},
		{CoapCodeServiceUnavailable, "503 Service Unavailable"},
		{CoapCodeGatewayTimeout, "504 Gateway Timeout"},
		{CoapCodeProxyingNotSupported, "505 Proxying Not Supported"},
		{CoapCode(255), "Unknown"},
	}

	for _, td := range testData {
		assert.Equal(t, td.codeString, CoapCodeToString(td.coapCode))
	}
}

func TestRouteMatching(t *testing.T) {

}

func TestMediaTypeUtils(t *testing.T) {
	assert.True(t, ValidCoapMediaTypeCode(MediaTypeTextPlain))
	assert.True(t, ValidCoapMediaTypeCode(MediaTypeOpaqueVndOmaLwm2m))

	assert.False(t, ValidCoapMediaTypeCode(MediaType(9999)))
}
