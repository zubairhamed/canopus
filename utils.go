package canopus

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

// Returns the string value for a Message Payload
func PayloadAsString(p MessagePayload) string {
	if p == nil {
		return ""
	}
	return p.String()
}

// GenerateMessageId generate a uint16 Message ID
func GenerateMessageID() uint16 {
	if CurrentMessageID != 65535 {
		CurrentMessageID++
	} else {
		CurrentMessageID = 1
	}
	return uint16(CurrentMessageID)
}

var genChars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

// GenerateToken generates a random token by a given length
func GenerateToken(l int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	token := make([]rune, l)
	for i := range token {
		token[i] = genChars[rand.Intn(len(genChars))]
	}
	return string(token)
}

// CoreResourcesFromString Converts to CoRE Resources Object from a CoRE String
func CoreResourcesFromString(str string) []*CoreResource {
	var re = regexp.MustCompile(`(<[^>]+>\s*(;\s*\w+\s*(=\s*(\w+|"([^"\\]*(\\.[^"\\]*)*)")\s*)?)*)`)
	var elemRe = regexp.MustCompile(`<[^>]*>`)

	var resources []*CoreResource
	m := re.FindAllString(str, -1)

	for _, match := range m {
		elemMatch := elemRe.FindString(match)
		target := elemMatch[1 : len(elemMatch)-1]

		resource := NewCoreResource()
		resource.Target = target

		if len(match) > len(elemMatch) {
			attrs := strings.Split(match[len(elemMatch)+1:], ";")

			for _, attr := range attrs {
				pair := strings.Split(attr, "=")

				resource.AddAttribute(pair[0], strings.Replace(pair[1], "\"", "", -1))
			}
		}
		resources = append(resources, resource)
	}
	return resources
}

// CoapCodeToString returns the string representation of a CoapCode
func CoapCodeToString(code CoapCode) string {
	switch code {
	case Get:
		return "GET"

	case Post:
		return "POST"

	case Put:
		return "PUT"

	case Delete:
		return "DELETE"

	case CoapCodeEmpty:
		return "0 Empty"

	case CoapCodeCreated:
		return "201 Created"

	case CoapCodeDeleted:
		return "202 Deleted"

	case CoapCodeValid:
		return "203 Valid"

	case CoapCodeChanged:
		return "204 Changed"

	case CoapCodeContent:
		return "205 Content"

	case CoapCodeBadRequest:
		return "400 Bad Request"

	case CoapCodeUnauthorized:
		return "401 Unauthorized"

	case CoapCodeBadOption:
		return "402 Bad Option"

	case CoapCodeForbidden:
		return "403 Forbidden"

	case CoapCodeNotFound:
		return "404 Not Found"

	case CoapCodeMethodNotAllowed:
		return "405 Method Not Allowed"

	case CoapCodeNotAcceptable:
		return "406 Not Acceptable"

	case CoapCodePreconditionFailed:
		return "412 Precondition Failed"

	case CoapCodeRequestEntityTooLarge:
		return "413 Request Entity Too Large"

	case CoapCodeUnsupportedContentFormat:
		return "415 Unsupported Content Format"

	case CoapCodeInternalServerError:
		return "500 Internal Server Error"

	case CoapCodeNotImplemented:
		return "501 Not Implemented"

	case CoapCodeBadGateway:
		return "502 Bad Gateway"

	case CoapCodeServiceUnavailable:
		return "503 Service Unavailable"

	case CoapCodeGatewayTimeout:
		return "504 Gateway Timeout"

	case CoapCodeProxyingNotSupported:
		return "505 Proxying Not Supported"

	default:
		return "Unknown"
	}
}

// ValidCoapMediaTypeCode Checks if a MediaType is of a valid code
func ValidCoapMediaTypeCode(mt MediaType) bool {
	switch mt {
	case MediaTypeTextPlain, MediaTypeTextXML, MediaTypeTextCsv, MediaTypeTextHTML, MediaTypeImageGif,
		MediaTypeImageJpeg, MediaTypeImagePng, MediaTypeImageTiff, MediaTypeAudioRaw, MediaTypeVideoRaw,
		MediaTypeApplicationLinkFormat, MediaTypeApplicationXML, MediaTypeApplicationOctetStream, MediaTypeApplicationRdfXML,
		MediaTypeApplicationSoapXML, MediaTypeApplicationAtomXML, MediaTypeApplicationXmppXML, MediaTypeApplicationExi,
		MediaTypeApplicationFastInfoSet, MediaTypeApplicationSoapFastInfoSet, MediaTypeApplicationJSON,
		MediaTypeApplicationXObitBinary, MediaTypeTextPlainVndOmaLwm2m, MediaTypeTlvVndOmaLwm2m,
		MediaTypeJSONVndOmaLwm2m, MediaTypeOpaqueVndOmaLwm2m:
		return true
	}

	return false
}

func logMsg(a ...interface{}) (n int, err error) {
	return fmt.Println(a)
}
