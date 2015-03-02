package goap
import "errors"

const (
    TYPE_CONFIRMABLE        = 0
    TYPE_NONCONFIRMABLE     = 1
    TYPE_ACKNOWLEDGEMENT    = 2
    TYPE_RESET              = 3
)

const (
    GET      = 1
    POST     = 2
    PUT      = 3
    DELETE   = 4
)

const (
    DATA_HEADER         = 0
    DATA_CODE           = 1
    DATA_MSGID_START    = 2
    DATA_MSGID_END      = 4
    DATA_TOKEN_START    = 4
)

const (
    OPTION_IF_MATCH         = 1
    OPTION_URI_HOST         = 3
    OPTION_ETAG             = 4
    OPTION_IF_NONE_MATCH    = 5
    OPTION_URI_PORT         = 7
    OPTION_LOCATION_PATH    = 8
    OPTION_URI_PATH         = 11
    OPTION_CONTENT_FORMAT   = 12
    OPTION_MAX_AGE          = 14
    OPTION_URI_QUERY        = 15
    OPTION_ACCEPT           = 17
    OPTION_LOCATION_QUERY   = 20
    OPTION_PROXY_URI        = 35
    OPTION_PROXY_SCHEME     = 39
    OPTION_SIZE1            = 60
)

const (
    MEDIATYPE_TEXT_PLAIN                = 0
    MEDIATYPE_APPLICATION_LINK_FORMAT   = 40
    MEDIATYPE_APPLICATION_XML           = 41
    MEDIATYPE_APPLICATION_OCTET_STREAM  = 42
    MEDIATYPE_APPLICATION_EXI           = 47
    MEDIATYPE_APPLICATION_JSON          = 50
)

const (
    COAPCODE_201_CREATED 					= 65
    COAPCODE_202_DELETED					= 66
    COAPCODE_203_VALID						= 67
    COAPCODE_204_CHANGED					= 68
    COAPCODE_205_CONTENT					= 69
    COAPCODE_400_BAD_REQUEST				= 128
    COAPCODE_401_UNAUTHORIZED				= 129
    COAPCODE_402_BAD_OPTION					= 130
    COAPCODE_403_FORBIDDEN					= 131
    COAPCODE_404_NOT_FOUND					= 132
    COAPCODE_405_METHOD_NOT_ALLOWED			= 133
    COAPCODE_406_NOT_ACCEPTABLE				= 134
    COAPCODE_412_PRECONDITION_FAILED		= 140
    COAPCODE_413_REQUEST_ENTITY_TOO_LARGE	= 141
    COAPCODE_415_UNSUPPORTED_CONTENT_FORMAT	= 143
    COAPCODE_500_INTERNAL_SERVER_ERROR		= 160
    COAPCODE_501_NOT_IMPLEMENTED			= 161
    COAPCODE_502_BAD_GATEWAY				= 162
    COAPCODE_503_SERVICE_UNAVAILABLE		= 163
    COAPCODE_504_GATEWAY_TIMEOUT			= 164
    COAPCODE_505_PROXYING_NOT_SUPPORTED		= 165
)

const PAYLOAD_MARKER = 0xff

const BUF_SIZE = 1500

const COAP_DEFAULT_HOST = ":5683"

// ERRORS
var ERR_NO_MATCHING_ROUTE = errors.New("No matching route found")
// return msg, errors.New("Packet length less than 4 bytes")
// return nil, errors.New("Invalid version")
// return msg, errors.New("Message Format Error. Option length has reserved value 15")
// return errors.New("Unknown message type")
// return errors.New("Invalid Token Length ( > 8)")
