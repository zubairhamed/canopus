package canopus
import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestGenerateMessageId(t *testing.T) {

	var id, id2 uint16
	id = GenerateMessageId()
	for i:=0; i < 100; i++ {
		id2 = id+1
		id = GenerateMessageId()
		assert.NotEqual(t, 65535, id)
		assert.Equal(t, id2, id)
	}
}

func TestGenerateToken(t *testing.T) {
	assert.Equal(t, "", GenerateToken(0))

	for i:=1; i < 10; i++ {
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

	assert.NotNil(t,resource1.GetAttribute("ct"))
	assert.Equal(t, "ct", resource1.GetAttribute("ct").Key)
	assert.Equal(t, "41", resource1.GetAttribute("ct").Value)

	assert.NotNil(t,resource1.GetAttribute("rt"))
	assert.Equal(t, "rt", resource1.GetAttribute("rt").Key)
	assert.Equal(t, "temperature-c", resource1.GetAttribute("rt").Value)

	assert.NotNil(t,resource1.GetAttribute("if"))
	assert.Equal(t, "if", resource1.GetAttribute("if").Key)
	assert.Equal(t, "sensor", resource1.GetAttribute("if").Value)

	resource2 := resources[1]
	assert.Equal(t, "/sensors/light", resource2.Target)
	assert.Equal(t, 3, len(resource2.Attributes))

	assert.NotNil(t,resource2.GetAttribute("ct"))
	assert.Equal(t, "ct", resource2.GetAttribute("ct").Key)
	assert.Equal(t, "41", resource2.GetAttribute("ct").Value)

	assert.NotNil(t,resource2.GetAttribute("rt"))
	assert.Equal(t, "rt", resource2.GetAttribute("rt").Key)
	assert.Equal(t, "light-lux", resource2.GetAttribute("rt").Value)

	assert.NotNil(t,resource2.GetAttribute("if"))
	assert.Equal(t, "if", resource2.GetAttribute("if").Key)
	assert.Equal(t, "sensor", resource2.GetAttribute("if").Value)
}

func TestCoapCodeToString(t *testing.T) {
	test_data := []struct {
		coapCode 	CoapCode
		codeString 	string
	}{
		{ GET, "GET" },
		{ POST, "POST" },
		{ PUT, "PUT" },
		{ DELETE, "DELETE" },
		{ COAPCODE_0_EMPTY, "0 Empty" },
		{ COAPCODE_201_CREATED, "201 Created" },
		{ COAPCODE_202_DELETED, "202 Deleted" },
		{ COAPCODE_203_VALID, "203 Valid" },
		{ COAPCODE_204_CHANGED, "204 Changed" },
		{ COAPCODE_205_CONTENT, "205 Content" },
		{ COAPCODE_400_BAD_REQUEST, "400 Bad Request" },
		{ COAPCODE_401_UNAUTHORIZED, "401 Unauthorized" },
		{ COAPCODE_402_BAD_OPTION, "402 Bad Option" },
		{ COAPCODE_403_FORBIDDEN, "403 Forbidden" },
		{ COAPCODE_404_NOT_FOUND, "404 Not Found" },
		{ COAPCODE_405_METHOD_NOT_ALLOWED, "405 Method Not Allowed" },
		{ COAPCODE_406_NOT_ACCEPTABLE, "406 Not Acceptable" },
		{ COAPCODE_412_PRECONDITION_FAILED, "412 Precondition Failed" },
		{ COAPCODE_413_REQUEST_ENTITY_TOO_LARGE, "413 Request Entity Too Large" },
		{ COAPCODE_415_UNSUPPORTED_CONTENT_FORMAT, "415 Unsupported Content Format" },
		{ COAPCODE_500_INTERNAL_SERVER_ERROR, "500 Internal Server Error" },
		{ COAPCODE_501_NOT_IMPLEMENTED, "501 Not Implemented" },
		{ COAPCODE_502_BAD_GATEWAY, "502 Bad Gateway" },
		{ COAPCODE_503_SERVICE_UNAVAILABLE, "503 Service Unavailable" },
		{ COAPCODE_504_GATEWAY_TIMEOUT, "504 Gateway Timeout" },
		{ COAPCODE_505_PROXYING_NOT_SUPPORTED, "505 Proxying Not Supported" },
		{ CoapCode(255), "Unknown" },
	}

	for _, td := range test_data {
		assert.Equal(t, td.codeString, CoapCodeToString(td.coapCode))
	}
}

func TestRouteMatching(t *testing.T) {

}

func TestMediaTypeUtils(t *testing.T) {
	assert.True(t, ValidCoapMediaTypeCode(MEDIATYPE_TEXT_PLAIN))
	assert.True(t, ValidCoapMediaTypeCode(MEDIATYPE_OPAQUE_VND_OMA_LWM2M))

	assert.False(t, ValidCoapMediaTypeCode(MediaType(9999)))
}
