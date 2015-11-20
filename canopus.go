package canopus

import (
	"errors"
	"math/rand"
	"time"
	"net"
)

// Message ID Generator, global
var MESSAGEID_CURR = 0

func init() {
	rand.Seed(time.Now().UTC().UnixNano())

	MESSAGEID_CURR = rand.Intn(65535)
}

// Types of Messages
const (
	TYPE_CONFIRMABLE     = 0
	TYPE_NONCONFIRMABLE  = 1
	TYPE_ACKNOWLEDGEMENT = 2
	TYPE_RESET           = 3
)

// Fragments/parts of a CoAP Message packet
const (
	DATA_HEADER      = 0
	DATA_CODE        = 1
	DATA_MSGID_START = 2
	DATA_MSGID_END   = 4
	DATA_TOKEN_START = 4
)

// OptionCode type represents a valid CoAP Option Code
type OptionCode int

const (
	OPTION_IF_MATCH       OptionCode = 1
	OPTION_URI_HOST       OptionCode = 3
	OPTION_ETAG           OptionCode = 4
	OPTION_IF_NONE_MATCH  OptionCode = 5
	OPTION_OBSERVE        OptionCode = 6
	OPTION_URI_PORT       OptionCode = 7
	OPTION_LOCATION_PATH  OptionCode = 8
	OPTION_URI_PATH       OptionCode = 11
	OPTION_CONTENT_FORMAT OptionCode = 12
	OPTION_MAX_AGE        OptionCode = 14
	OPTION_URI_QUERY      OptionCode = 15
	OPTION_ACCEPT         OptionCode = 17
	OPTION_LOCATION_QUERY OptionCode = 20
	OPTION_BLOCK2         OptionCode = 23
	OPTION_BLOCK1         OptionCode = 27
	OPTION_SIZE2          OptionCode = 28
	OPTION_PROXY_URI      OptionCode = 35
	OPTION_PROXY_SCHEME   OptionCode = 39
	OPTION_SIZE1          OptionCode = 60
)

// Valid CoAP Codes
type CoapCode uint8

const (
	GET    CoapCode = 1
	POST   CoapCode = 2
	PUT    CoapCode = 3
	DELETE CoapCode = 4

	COAPCODE_0_EMPTY                        CoapCode = 0
	COAPCODE_201_CREATED                    CoapCode = 65
	COAPCODE_202_DELETED                    CoapCode = 66
	COAPCODE_203_VALID                      CoapCode = 67
	COAPCODE_204_CHANGED                    CoapCode = 68
	COAPCODE_205_CONTENT                    CoapCode = 69
	COAPCODE_400_BAD_REQUEST                CoapCode = 128
	COAPCODE_401_UNAUTHORIZED               CoapCode = 129
	COAPCODE_402_BAD_OPTION                 CoapCode = 130
	COAPCODE_403_FORBIDDEN                  CoapCode = 131
	COAPCODE_404_NOT_FOUND                  CoapCode = 132
	COAPCODE_405_METHOD_NOT_ALLOWED         CoapCode = 133
	COAPCODE_406_NOT_ACCEPTABLE             CoapCode = 134
	COAPCODE_409_CONFLICT                   CoapCode = 137
	COAPCODE_412_PRECONDITION_FAILED        CoapCode = 140
	COAPCODE_413_REQUEST_ENTITY_TOO_LARGE   CoapCode = 141
	COAPCODE_415_UNSUPPORTED_CONTENT_FORMAT CoapCode = 143
	COAPCODE_500_INTERNAL_SERVER_ERROR      CoapCode = 160
	COAPCODE_501_NOT_IMPLEMENTED            CoapCode = 161
	COAPCODE_502_BAD_GATEWAY                CoapCode = 162
	COAPCODE_503_SERVICE_UNAVAILABLE        CoapCode = 163
	COAPCODE_504_GATEWAY_TIMEOUT            CoapCode = 164
	COAPCODE_505_PROXYING_NOT_SUPPORTED     CoapCode = 165
)

// Default Acknowledgement Timeout
const DEFAULT_ACK_TIMEOUT = 2
const DEFAULT_ACK_RANDOM_FACTOR = 1.5
const DEFAULT_MAX_RETRANSMIT = 4
const DEFAULT_NSTART = 1
const DEFAULT_LEISURE = 5
const DEFAULT_PROBING_RATE = 1

const COAP_DEFAULT_HOST = ""
const COAP_DEFAULT_PORT = 5683
const COAPS_DEFAULT_PORT = 5684

const PAYLOAD_MARKER = 0xff
const MAX_PACKET_SIZE = 1500

// Number of seconds before a MessageID Purge is initiated
const MESSAGEID_PURGE_DURATION = 60

type RouteHandler func(CoapRequest) CoapResponse

// type ResponseHandler func(CoapRespose, error)

// Supported Media Types
type MediaType int

const (
	MEDIATYPE_TEXT_PLAIN                  MediaType = 0
	MEDIATYPE_TEXT_XML                    MediaType = 1
	MEDIATYPE_TEXT_CSV                    MediaType = 2
	MEDIATYPE_TEXT_HTML                   MediaType = 3
	MEDIATYPE_IMAGE_GIF                   MediaType = 21
	MEDIATYPE_IMAGE_JPEG                  MediaType = 22
	MEDIATYPE_IMAGE_PNG                   MediaType = 23
	MEDIATYPE_IMAGE_TIFF                  MediaType = 24
	MEDIATYPE_AUDIO_RAW                   MediaType = 25
	MEDIATYPE_VIDEO_RAW                   MediaType = 26
	MEDIATYPE_APPLICATION_LINK_FORMAT     MediaType = 40
	MEDIATYPE_APPLICATION_XML             MediaType = 41
	MEDIATYPE_APPLICATION_OCTET_STREAM    MediaType = 42
	MEDIATYPE_APPLICATION_RDFXML          MediaType = 43
	MEDIATYPE_APPLICATION_SOAPXML         MediaType = 44
	MEDIATYPE_APPLICATION_ATOMXML         MediaType = 45
	MEDIATYPE_APPLICATION_XMPPXML         MediaType = 46
	MEDIATYPE_APPLICATION_EXI             MediaType = 47
	MEDIATYPE_APPLICATION_FASTINFOSET     MediaType = 48
	MEDIATYPE_APPLICATION_SOAPFASTINFOSET MediaType = 49
	MEDIATYPE_APPLICATION_JSON            MediaType = 50
	MEDIATYPE_APPLICATION_X_OBIT_BINARY   MediaType = 51
	MEDIATYPE_TEXT_PLAIN_VND_OMA_LWM2M    MediaType = 1541
	MEDIATYPE_TLV_VND_OMA_LWM2M           MediaType = 1542
	MEDIATYPE_JSON_VND_OMA_LWM2M          MediaType = 1543
	MEDIATYPE_OPAQUE_VND_OMA_LWM2M        MediaType = 1544
)

const (
	METHOD_GET     = "GET"
	METHOD_PUT     = "PUT"
	METHOD_POST    = "POST"
	METHOD_DELETE  = "DELETE"
	METHOD_OPTIONS = "OPTIONS"
	METHOD_PATCH   = "PATCH"
)

// ERRORS
var ERR_PACKET_LENGTH_LESS_THAN_4 = errors.New("Packet length less than 4 bytes")
var ERR_INVALID_VERSION = errors.New("Invalid CoAP version. Should be 1.")
var ERR_OPTION_LENGTH_USES_VALUE_15 = errors.New(("Message format error. Option length has reserved value of 15"))
var ERR_OPTION_DELTA_USES_VALUE_15 = errors.New(("Message format error. Option delta has reserved value of 15"))
var ERR_UNKNOWN_MESSAGE_TYPE = errors.New("Unknown message type")
var ERR_INVALID_TOKEN_LENGTH = errors.New("Invalid Token Length ( > 8)")
var ERR_UNKNOWN_CRITICAL_OPTION = errors.New("Unknown critical option encountered")
var ERR_UNSUPPORTED_METHOD = errors.New("Unsupported Method")
var ERR_NO_MATCHING_ROUTE = errors.New("No matching route found")
var ERR_UNSUPPORTED_CONTENT_FORMAT = errors.New("Unsupported Content-Format")
var ERR_NO_MATCHING_METHOD = errors.New("No matching method")
var ERR_NIL_MESSAGE = errors.New("Message is nil")
var ERR_NIL_CONN = errors.New("Connection object is nil")
var ERR_NIL_ADDR = errors.New("Address cannot be nil")

//// API ////
type CoapServer interface {
	Start()
	Stop()
	SetProxyFilter(fn ProxyFilter)
	Get(path string, fn RouteHandler) *Route
	Delete(path string, fn RouteHandler) *Route
	Put(path string, fn RouteHandler) *Route
	Post(path string, fn RouteHandler) *Route
	Options(path string, fn RouteHandler) *Route
	Patch(path string, fn RouteHandler) *Route
	NewRoute(path string, method CoapCode, fn RouteHandler) *Route
	Send(req CoapRequest) (CoapResponse, error)
	SendTo(req CoapRequest, addr *net.UDPAddr) (CoapResponse, error)
	NotifyChange(resource, value string, confirm bool)
	Dial(host string)
	Dial6(host string)
	OnNotify(fn FnEventNotify)
	OnStart(fn FnEventStart)
	OnClose(fn FnEventClose)
	OnDiscover(fn FnEventDiscover)
	OnError(fn FnEventError)
	OnObserve(fn FnEventObserve)
	OnObserveCancel(fn FnEventObserveCancel)
	OnMessage(fn FnEventMessage)
	ProxyHttp(enabled bool)
	ProxyCoap(enabled bool)
	GetEvents() *CanopusEvents
	GetLocalAddress() *net.UDPAddr

	AllowProxyForwarding(*Message, *net.UDPAddr) (bool)
	GetRoutes() []*Route
	ForwardCoap(msg *Message, conn *net.UDPConn, addr *net.UDPAddr)
	ForwardHttp(msg *Message, conn *net.UDPConn, addr *net.UDPAddr)

	AddObservation(resource, token string, addr *net.UDPAddr)
	HasObservation(resource string, addr *net.UDPAddr) bool
	RemoveObservation(resource string, addr *net.UDPAddr)

	IsDuplicateMessage(msg *Message) bool
	UpdateMessageTS(msg *Message)
}