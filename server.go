package canopus

import (
	"bytes"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

type ProxyType int

const (
	ProxyHTTP ProxyType = 0
	ProxyCOAP ProxyType = 1
)

func NewLocalServer() CoapServer {
	return NewServer("5683", "")
}

func NewCoapServer(local string) CoapServer {
	return NewServer(local, "")
}

func NewCoapClient() CoapServer {
	return NewServer("0", "")
}

func NewServer(local, remote string) CoapServer {
	localHost := local
	if !strings.Contains(localHost, ":") {
		localHost = ":" + localHost
	}
	localAddr, _ := net.ResolveUDPAddr("udp6", localHost)

	var remoteAddr *net.UDPAddr
	if remote != "" {
		remoteHost := remote
		if !strings.Contains(remoteHost, ":") {
			remoteHost = ":" + remoteHost
		}
		remoteAddr, _ = net.ResolveUDPAddr("udp6", remoteHost)
	}

	return &DefaultCoapServer{
		remoteAddr:        remoteAddr,
		localAddr:         localAddr,
		events:            NewEvents(),
		observations:      make(map[string][]*Observation),
		fnHandleCOAPProxy: NullProxyHandler,
		fnHandleHTTPProxy: NullProxyHandler,
		fnProxyFilter:     NullProxyFilter,
		stopChannel:       make(chan int),
	}
}

type DefaultCoapServer struct {
	localAddr  *net.UDPAddr
	remoteAddr *net.UDPAddr

	localConn  *net.UDPConn
	remoteConn *net.UDPConn

	messageIds    map[uint16]time.Time
	blockMessages map[string][]*BlockMessage

	routes       []*Route
	events       *Events
	observations map[string][]*Observation

	fnHandleHTTPProxy ProxyHandler
	fnHandleCOAPProxy ProxyHandler
	fnProxyFilter     ProxyFilter

	stopChannel chan int
}

func (s *DefaultCoapServer) GetEvents() *Events {
	return s.events
}

func (s *DefaultCoapServer) Start() {
	var discoveryRoute RouteHandler = func(req CoapRequest) CoapResponse {
		msg := req.GetMessage()

		ack := ContentMessage(msg.MessageID, MessageAcknowledgment)
		ack.Token = make([]byte, len(msg.Token))
		copy(ack.Token, msg.Token)

		ack.AddOption(OptionContentFormat, MediaTypeApplicationLinkFormat)

		var buf bytes.Buffer
		for _, r := range s.routes {
			if r.Path != ".well-known/core" {
				buf.WriteString("</")
				buf.WriteString(r.Path)
				buf.WriteString(">")

				// Media Types
				lenMt := len(r.MediaTypes)
				if lenMt > 0 {
					buf.WriteString(";ct=")
					for idx, mt := range r.MediaTypes {

						buf.WriteString(strconv.Itoa(int(mt)))
						if idx+1 < lenMt {
							buf.WriteString(" ")
						}
					}
				}

				buf.WriteString(",")
				// buf.WriteString("</" + r.Path + ">;ct=0,")
			}
		}
		ack.Payload = NewPlainTextPayload(buf.String())

		resp := NewResponseWithMessage(ack)

		return resp
	}

	s.NewRoute("/.well-known/core", Get, discoveryRoute)
	s.serveServer()
}

func (s *DefaultCoapServer) serveServer() {
	s.messageIds = make(map[uint16]time.Time)
	s.blockMessages = make(map[string][]*BlockMessage)

	conn, err := net.ListenUDP("udp", s.localAddr)
	if err != nil {
		s.events.Error(err)
		log.Fatal(err)
	}

	s.localConn = conn

	if conn == nil {
		log.Fatal("An error occured starting up CoAP Server")
	} else {
		log.Println("Started CoAP Server ", conn.LocalAddr())
	}

	s.events.Started(s)
	s.handleMessageIDPurge()

	readBuf := make([]byte, MaxPacketSize)
	for {
		select {
		case <-s.stopChannel:
			return

		default:
			// continue
		}

		len, addr, err := conn.ReadFromUDP(readBuf)

		if err == nil {

			msgBuf := make([]byte, len)
			copy(msgBuf, readBuf)

			go s.handleMessage(msgBuf, conn, addr)
		}
	}
}

func (s *DefaultCoapServer) Stop() {
	s.localConn.Close()
	close(s.stopChannel)
}


func (s *DefaultCoapServer) UpdateBlockMessage(client string, msg *Message, seq uint32) {
	msgs := s.blockMessages[client]
	if msgs == nil {
		msgs = []*BlockMessage{}
	}

	bm := NewBlockMessage()
	bm.StoredMessage = msg
	bm.Sequence = seq
	s.blockMessages[client] = append(msgs, bm)
}

//func (s *DefaultCoapServer) GetBlockBuffer(client string) []byte {
//	return s.blockMessages[client]
//}
//
//func (s *DefaultCoapServer) PurgeBlockBuffer(client string) {
//	s.blockMessages[client] = nil
//}

func (s *DefaultCoapServer) handleMessageIDPurge() {
	// Routine for clearing up message IDs which has expired
	ticker := time.NewTicker(MessageIDPurgeDuration * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				for k, v := range s.messageIds {
					elapsed := time.Since(v)
					if elapsed > MessageIDPurgeDuration {
						delete(s.messageIds, k)
					}
				}
			}
		}
	}()
}

func (s *DefaultCoapServer) SetProxyFilter(fn ProxyFilter) {
	s.fnProxyFilter = fn
}

func (s *DefaultCoapServer) handleMessage(msgBuf []byte, conn *net.UDPConn, addr *net.UDPAddr) {
	msg, err := BytesToMessage(msgBuf)
	s.events.Message(msg, true)

	if msg.MessageType == MessageAcknowledgment {
		handleResponse(s, msg, conn, addr)
	} else {
		handleRequest(s, err, msg, conn, addr)
	}
}

func (s *DefaultCoapServer) Get(path string, fn RouteHandler) *Route {
	return s.add(MethodGet, path, fn)
}

func (s *DefaultCoapServer) Delete(path string, fn RouteHandler) *Route {
	return s.add(MethodDelete, path, fn)
}

func (s *DefaultCoapServer) Put(path string, fn RouteHandler) *Route {
	return s.add(MethodPut, path, fn)
}

func (s *DefaultCoapServer) Post(path string, fn RouteHandler) *Route {
	return s.add(MethodPost, path, fn)
}

func (s *DefaultCoapServer) Options(path string, fn RouteHandler) *Route {
	return s.add(MethodOptions, path, fn)
}

func (s *DefaultCoapServer) Patch(path string, fn RouteHandler) *Route {
	return s.add(MethodPatch, path, fn)
}

func (s *DefaultCoapServer) add(method string, path string, fn RouteHandler) *Route {
	route := CreateNewRoute(path, method, fn)
	s.routes = append(s.routes, route)

	return route
}

func (s *DefaultCoapServer) NewRoute(path string, method CoapCode, fn RouteHandler) *Route {
	route := CreateNewRoute(path, MethodString(method), fn)
	s.routes = append(s.routes, route)

	return route
}

func (s *DefaultCoapServer) Send(req CoapRequest) (CoapResponse, error) {
	s.events.Message(req.GetMessage(), false)
	response, err := SendMessageTo(req.GetMessage(), NewUDPConnection(s.localConn), s.remoteAddr)

	if err != nil {
		s.events.Error(err)
		return response, err
	}
	s.events.Message(response.GetMessage(), true)

	return response, err
}

func (s *DefaultCoapServer) SendTo(req CoapRequest, addr *net.UDPAddr) (CoapResponse, error) {
	return SendMessageTo(req.GetMessage(), NewUDPConnection(s.localConn), addr)
}

func (s *DefaultCoapServer) NotifyChange(resource, value string, confirm bool) {
	t := s.observations[resource]

	if t != nil {
		var req CoapRequest

		if confirm {
			req = NewRequest(MessageConfirmable, CoapCodeContent, GenerateMessageID())
		} else {
			req = NewRequest(MessageAcknowledgment, CoapCodeContent, GenerateMessageID())
		}

		for _, r := range t {
			req.SetToken(r.Token)
			req.SetStringPayload(value)
			req.SetRequestURI(r.Resource)
			r.NotifyCount++
			req.GetMessage().AddOption(OptionObserve, r.NotifyCount)

			go s.SendTo(req, r.Addr)
		}
	}
}

func (s *DefaultCoapServer) AddObservation(resource, token string, addr *net.UDPAddr) {
	s.observations[resource] = append(s.observations[resource], NewObservation(addr, token, resource))
}

func (s *DefaultCoapServer) HasObservation(resource string, addr *net.UDPAddr) bool {
	obs := s.observations[resource]
	if obs == nil {
		return false
	}

	for _, o := range obs {
		if o.Addr.String() == addr.String() {
			return true
		}
	}
	return false
}

func (s *DefaultCoapServer) RemoveObservation(resource string, addr *net.UDPAddr) {
	obs := s.observations[resource]
	if obs == nil {
		return
	}

	for idx, o := range obs {
		if o.Addr.String() == addr.String() {
			s.observations[resource] = append(obs[:idx], obs[idx+1:]...)
			return
		}
	}
}

func (s *DefaultCoapServer) Dial(host string) {
	s.Dial6(host)
}

func (s *DefaultCoapServer) Dial6(host string) {
	remoteAddr, _ := net.ResolveUDPAddr("udp6", host)

	s.remoteAddr = remoteAddr
}

func (s *DefaultCoapServer) OnNotify(fn FnEventNotify) {
	s.events.OnNotify(fn)
}

func (s *DefaultCoapServer) OnStart(fn FnEventStart) {
	s.events.OnStart(fn)
}

func (s *DefaultCoapServer) OnClose(fn FnEventClose) {
	s.events.OnClose(fn)
}

func (s *DefaultCoapServer) OnDiscover(fn FnEventDiscover) {
	s.events.OnDiscover(fn)
}

func (s *DefaultCoapServer) OnError(fn FnEventError) {
	s.events.OnError(fn)
}

func (s *DefaultCoapServer) OnObserve(fn FnEventObserve) {
	s.events.OnObserve(fn)
}

func (s *DefaultCoapServer) OnObserveCancel(fn FnEventObserveCancel) {
	s.events.OnObserveCancel(fn)
}

func (s *DefaultCoapServer) OnMessage(fn FnEventMessage) {
	s.events.OnMessage(fn)
}

func (s *DefaultCoapServer) OnBlockMessage(fn FnEventBlockMessage) {
	s.events.OnBlockMessage(fn)
}

func (s *DefaultCoapServer) ProxyHTTP(enabled bool) {
	if enabled {
		s.fnHandleHTTPProxy = HTTPProxyHandler
	} else {
		s.fnHandleHTTPProxy = NullProxyHandler
	}
}

func (s *DefaultCoapServer) ProxyCoap(enabled bool) {
	if enabled {
		s.fnHandleCOAPProxy = COAPProxyHandler
	} else {
		s.fnHandleCOAPProxy = NullProxyHandler
	}
}

func (s *DefaultCoapServer) AllowProxyForwarding(msg *Message, addr *net.UDPAddr) bool {
	return s.fnProxyFilter(msg, addr)
}

func (s *DefaultCoapServer) ForwardCoap(msg *Message, conn *net.UDPConn, addr *net.UDPAddr) {
	s.fnHandleCOAPProxy(msg, conn, addr)
}

func (s *DefaultCoapServer) ForwardHTTP(msg *Message, conn *net.UDPConn, addr *net.UDPAddr) {
	s.fnHandleHTTPProxy(msg, conn, addr)
}

func (s *DefaultCoapServer) GetRoutes() []*Route {
	return s.routes
}

func (s *DefaultCoapServer) GetLocalAddress() *net.UDPAddr {
	return s.localAddr
}

func (s *DefaultCoapServer) IsDuplicateMessage(msg *Message) bool {
	_, ok := s.messageIds[msg.MessageID]

	return ok
}

func (s *DefaultCoapServer) UpdateMessageTS(msg *Message) {
	s.messageIds[msg.MessageID] = time.Now()
}

func NewObservation(addr *net.UDPAddr, token string, resource string) *Observation {
	return &Observation{
		Addr:        addr,
		Token:       token,
		Resource:    resource,
		NotifyCount: 0,
	}
}

type Observation struct {
	Addr        *net.UDPAddr
	Token       string
	Resource    string
	NotifyCount int
}
