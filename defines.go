package goap
import "errors"

const (
    TYPE_CONFIRMABLE        = 0
    TYPE_NONCONFIRMABLE     = 1
    TYPE_ACKNOWLEDGEMENT    = 2
    TYPE_RESET              = 3
)

const (
    DATA_HEADER         = 0
    DATA_CODE           = 1
    DATA_MSGID_START    = 2
    DATA_MSGID_END      = 4
    DATA_TOKEN_START    = 4
)

type OptionCode int
const (
    OPTION_IF_MATCH         OptionCode = 1
    OPTION_URI_HOST         OptionCode = 3
    OPTION_ETAG             OptionCode = 4
    OPTION_IF_NONE_MATCH    OptionCode = 5
    OPTION_URI_PORT         OptionCode = 7
    OPTION_LOCATION_PATH    OptionCode = 8
    OPTION_URI_PATH         OptionCode = 11
    OPTION_CONTENT_FORMAT   OptionCode = 12
    OPTION_MAX_AGE          OptionCode = 14
    OPTION_URI_QUERY        OptionCode = 15
    OPTION_ACCEPT           OptionCode = 17
    OPTION_LOCATION_QUERY   OptionCode = 20
    OPTION_PROXY_URI        OptionCode = 35
    OPTION_PROXY_SCHEME     OptionCode = 39
    OPTION_SIZE1            OptionCode = 60
)

type MediaType byte
const (
    MEDIATYPE_TEXT_PLAIN               MediaType = 0
    MEDIATYPE_APPLICATION_LINK_FORMAT  MediaType = 40
    MEDIATYPE_APPLICATION_XML          MediaType = 41
    MEDIATYPE_APPLICATION_OCTET_STREAM MediaType = 42
    MEDIATYPE_APPLICATION_EXI          MediaType = 47
    MEDIATYPE_APPLICATION_JSON         MediaType = 50
)

type CoapCode uint8

const (
    GET      CoapCode = 1
    POST     CoapCode = 2
    PUT      CoapCode = 3
    DELETE   CoapCode = 4
)

const (
	COAPCODE_0_EMPTY						= 0
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

const DEFAULT_ACK_TIMEOUT 			= 2
const DEFAULT_ACK_RANDOM_FACTOR 	= 1.5
const DEFAULT_MAX_RETRANSMIT 		= 4
const DEFAULT_NSTART 				= 1
const DEFAULT_LEISURE 				= 5
const DEFAULT_PROBING_RATE 			= 1

/*
func MaxTransmitSpan() {
	return  ACK_TIMEOUT * ((2 ** MAX_RETRANSMIT) - 1) * ACK_RANDOM_FACTOR
}

func MaxTransmitWait() {
	return ACK_TIMEOUT * ((2 ** (MAX_RETRANSMIT + 1)) - 1) * ACK_RANDOM_FACTOR
*/

const COAP_DEFAULT_HOST     = ":5683"
const COAPS_DEFAULT_HOST    = ":5684"

const PAYLOAD_MARKER = 0xff
const BUF_SIZE = 1500

const MESSAGEID_PURGE_DURATION = 60

// ERRORS
var ERR_NO_MATCHING_ROUTE = errors.New("No matching route found")
var ERR_PACKET_LENGTH_LESS_THAN_4 = errors.New("Packet length less than 4 bytes")
var ERR_INVALID_VERSION = errors.New("Invalid CoAP version. Should be 1.")
var ERR_OPTION_LENGTH_USES_VALUE_15 = errors.New(("Message format error. Option length has reserved value of 15"))
var ERR_UNKNOWN_MESSAGE_TYPE = errors.New("Unknown message type")
var ERR_INVALID_TOKEN_LENGTH = errors.New("Invalid Token Length ( > 8)")
