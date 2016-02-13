package canopus

import (
	"errors"
	"math/rand"
	"net"
	"time"
)

// Message ID Generator, global
var CurrentMessageId = 0

func init() {
	rand.Seed(time.Now().UTC().UnixNano())

	CurrentMessageId = rand.Intn(65535)
}

// Types of Messages
const (
	MessageConfirmable     = 0
	MessageNonConfirmable  = 1
	MessageAcknowledgement = 2
	MessageReset           = 3
)

// Fragments/parts of a CoAP Message packet
const (
	DataHeader     = 0
	DataCode       = 1
	DataMsgIdStart = 2
	DataMsgIdEnd   = 4
	DataTokenStart = 4
)

// OptionCode type represents a valid CoAP Option Code
type OptionCode int

const (
	OptionIfMatch       OptionCode = 1
	OptionUriHost       OptionCode = 3
	OptionEtag          OptionCode = 4
	OptionIfNoneMatch   OptionCode = 5
	OptionObserve       OptionCode = 6
	OptionUriPort       OptionCode = 7
	OptionLocationPath  OptionCode = 8
	OptionUriPath       OptionCode = 11
	OptionContentFormat OptionCode = 12
	OptionMaxAge        OptionCode = 14
	OptionUriQuery      OptionCode = 15
	OptionAccept        OptionCode = 17
	OptionLocationQuery OptionCode = 20
	OptionBlock2        OptionCode = 23
	OptionBlock1        OptionCode = 27
	OptionSize2         OptionCode = 28
	OptionProxyUri      OptionCode = 35
	OptionProxyScheme   OptionCode = 39
	OptionSize1         OptionCode = 60
)

// Valid CoAP Codes
type CoapCode uint8

const (
	Get    CoapCode = 1
	Post   CoapCode = 2
	Put    CoapCode = 3
	Delete CoapCode = 4

	CoapCode_Empty                    CoapCode = 0
	CoapCode_Created                  CoapCode = 65
	CoapCode_Deleted                  CoapCode = 66
	CoapCode_Valid                    CoapCode = 67
	CoapCode_Changed                  CoapCode = 68
	CoapCode_Content                  CoapCode = 69
	CoapCode_BadRequest               CoapCode = 128
	CoapCode_Unauthorized             CoapCode = 129
	CoapCode_BadOption                CoapCode = 130
	CoapCode_Forbidden                CoapCode = 131
	CoapCode_NotFound                 CoapCode = 132
	CoapCode_MethodNotAllowed         CoapCode = 133
	CoapCode_NotAcceptable            CoapCode = 134
	CoapCode_Conflict                 CoapCode = 137
	CoapCode_PreconditionFailed       CoapCode = 140
	CoapCode_RequestEntityTooLarge    CoapCode = 141
	CoapCode_UnsupportedContentFormat CoapCode = 143
	CoapCode_InternalServerError      CoapCode = 160
	CoapCode_NotImplemented           CoapCode = 161
	CoapCode_BadGateway               CoapCode = 162
	CoapCode_ServiceUnavailable       CoapCode = 163
	CoapCode_GatewayTimeout           CoapCode = 164
	CoapCode_ProxyingNotSupported     CoapCode = 165
)

// Default Acknowledgement Timeout
const DefaultAckTimeout = 2
const DefaultAckRandomFactor = 1.5
const DefaultMaxRetransmit = 4
const DefaultNStart = 1
const DefaultLeisure = 5
const DefaultProbingRate = 1

const CoapDefaultHost = ""
const CoapDefaultPort = 5683
const CoapsDefaultPort = 5684

const PayloadMarker = 0xff
const MaxPacketSize = 1500

// Number of seconds before a MessageID Purge is initiated
const MessageIdPurgeDuration = 60

type RouteHandler func(CoapRequest) CoapResponse

// type ResponseHandler func(CoapRespose, error)

// Supported Media Types
type MediaType int

const (
	MediaTypeTextPlain                  MediaType = 0
	MediaTypeTextXml                    MediaType = 1
	MediaTypeTextCsv                    MediaType = 2
	MediaTypeTextHtml                   MediaType = 3
	MediaTypeImageGif                   MediaType = 21
	MediaTypeImageJpeg                  MediaType = 22
	MediaTypeImagePng                   MediaType = 23
	MediaTypeImageTiff                  MediaType = 24
	MediaTypeAudioRaw                   MediaType = 25
	MediaTypeVideoRaw                   MediaType = 26
	MediaTypeApplicationLinkFormat      MediaType = 40
	MediaTypeApplicationXml             MediaType = 41
	MediaTypeApplicationOctetStream     MediaType = 42
	MediaTypeApplicationRdfXml          MediaType = 43
	MediaTypeApplicationSoapXml         MediaType = 44
	MediaTypeApplicationAtomXml         MediaType = 45
	MediaTypeApplicationXmppXml         MediaType = 46
	MediaTypeApplicationExi             MediaType = 47
	MediaTypeApplicationFastInfoSet     MediaType = 48
	MediaTypeApplicationSoapFastInfoSet MediaType = 49
	MediaTypeApplicationJson            MediaType = 50
	MediaTypeApplicationXObitBinary     MediaType = 51
	MediaTypeTextPlainVndOmaLwm2m       MediaType = 1541
	MediaTypeTlvVndOmaLwm2m             MediaType = 1542
	MediaTypeJsonVndOmaLwm2m            MediaType = 1543
	MediaTypeOpaqueVndOmaLwm2m          MediaType = 1544
)

const (
	MethodGet      = "GET"
	MethodPut      = "PUT"
	METHOD_POST    = "POST"
	MethodDelete   = "DELETE"
	METHOD_OPTIONS = "OPTIONS"
	METHOD_PATCH   = "PATCH"
)

// ERRORS
var ErrPacketLengthLessThan4 = errors.New("Packet length less than 4 bytes")
var ErrInvalidCoapVersion = errors.New("Invalid CoAP version. Should be 1.")
var ErrOptionLengthUsesValue15 = errors.New(("Message format error. Option length has reserved value of 15"))
var ErrOptionDeltaUsesValue15 = errors.New(("Message format error. Option delta has reserved value of 15"))
var ErrUnknownMessageType = errors.New("Unknown message type")
var ErrInvalidTokenLength = errors.New("Invalid Token Length ( > 8)")
var ErrUnknownCriticalOption = errors.New("Unknown critical option encountered")
var ErrUnsupportedMethod = errors.New("Unsupported Method")
var ErrNoMatchingRoute = errors.New("No matching route found")
var ErrUnsupportedContentFormat = errors.New("Unsupported Content-Format")
var ErrNoMatchingMethod = errors.New("No matching method")
var ErrNilMessage = errors.New("Message is nil")
var ErrNilConn = errors.New("Connection object is nil")
var ErrNilAddr = errors.New("Address cannot be nil")

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

	AllowProxyForwarding(*Message, *net.UDPAddr) bool
	GetRoutes() []*Route
	ForwardCoap(msg *Message, conn *net.UDPConn, addr *net.UDPAddr)
	ForwardHttp(msg *Message, conn *net.UDPConn, addr *net.UDPAddr)

	AddObservation(resource, token string, addr *net.UDPAddr)
	HasObservation(resource string, addr *net.UDPAddr) bool
	RemoveObservation(resource string, addr *net.UDPAddr)

	IsDuplicateMessage(msg *Message) bool
	UpdateMessageTS(msg *Message)
}

// A simple wrapper interface around a connection
// This was primarily concieved so that mocks could be
// created to unit test connection code
type CanopusConnection interface {
	GetConnection() net.Conn
	Write(b []byte) (int, error)
	SetReadDeadline(t time.Time) error
	Read() (buf []byte, n int, err error)
	WriteTo(b []byte, addr net.Addr) (int, error)
}
