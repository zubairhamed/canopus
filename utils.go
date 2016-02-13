package canopus

import (
	"math/rand"
	"regexp"
	"strings"
	"time"
)

// GenerateMessageId generate a uint16 Message ID
func GenerateMessageId() uint16 {
	if CurrentMessageId != 65535 {
		CurrentMessageId++
	} else {
		CurrentMessageId = 1
	}
	return uint16(CurrentMessageId)
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

	case CoapCode_Empty:
		return "0 Empty"

	case CoapCode_Created:
		return "201 Created"

	case CoapCode_Deleted:
		return "202 Deleted"

	case CoapCode_Valid:
		return "203 Valid"

	case CoapCode_Changed:
		return "204 Changed"

	case CoapCode_Content:
		return "205 Content"

	case CoapCode_BadRequest:
		return "400 Bad Request"

	case CoapCode_Unauthorized:
		return "401 Unauthorized"

	case CoapCode_BadOption:
		return "402 Bad Option"

	case CoapCode_Forbidden:
		return "403 Forbidden"

	case CoapCode_NotFound:
		return "404 Not Found"

	case CoapCode_MethodNotAllowed:
		return "405 Method Not Allowed"

	case CoapCode_NotAcceptable:
		return "406 Not Acceptable"

	case CoapCode_PreconditionFailed:
		return "412 Precondition Failed"

	case CoapCode_RequestEntityTooLarge:
		return "413 Request Entity Too Large"

	case CoapCode_UnsupportedContentFormat:
		return "415 Unsupported Content Format"

	case CoapCode_InternalServerError:
		return "500 Internal Server Error"

	case CoapCode_NotImplemented:
		return "501 Not Implemented"

	case CoapCode_BadGateway:
		return "502 Bad Gateway"

	case CoapCode_ServiceUnavailable:
		return "503 Service Unavailable"

	case CoapCode_GatewayTimeout:
		return "504 Gateway Timeout"

	case CoapCode_ProxyingNotSupported:
		return "505 Proxying Not Supported"

	default:
		return "Unknown"
	}
}

// ValidCoapMediaTypeCode Checks if a MediaType is of a valid code
func ValidCoapMediaTypeCode(mt MediaType) bool {
	switch mt {
	case MediaTypeTextPlain, MediaTypeTextXml, MediaTypeTextCsv, MediaTypeTextHtml, MediaTypeImageGif,
		MediaTypeImageJpeg, MediaTypeImagePng, MediaTypeImageTiff, MediaTypeAudioRaw, MediaTypeVideoRaw,
		MediaTypeApplicationLinkFormat, MediaTypeApplicationXml, MediaTypeApplicationOctetStream, MediaTypeApplicationRdfXml,
		MediaTypeApplicationSoapXml, MediaTypeApplicationAtomXml, MediaTypeApplicationXmppXml, MediaTypeApplicationExi,
		MediaTypeApplicationFastInfoSet, MediaTypeApplicationSoapFastInfoSet, MediaTypeApplicationJson,
		MediaTypeApplicationXObitBinary, MediaTypeTextPlainVndOmaLwm2m, MediaTypeTlvVndOmaLwm2m,
		MediaTypeJsonVndOmaLwm2m, MediaTypeOpaqueVndOmaLwm2m:
		return true
	}

	return false
}
