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

type ServerConfiguration struct {
	EnableResourceDiscovery bool
}

func NewServer() CoapServer {
	return createServer()
}

func createServer() CoapServer {

	return &DefaultCoapServer{
		events:                  NewEvents(),
		observations:            make(map[string][]*Observation),
		fnHandleCOAPProxy:       NullProxyHandler,
		fnHandleHTTPProxy:       NullProxyHandler,
		fnProxyFilter:           NullProxyFilter,
		stopChannel:             make(chan int),
		coapResponseChannelsMap: make(map[uint16]chan *CoapResponseChannel),
		messageIds:              make(map[uint16]time.Time),
		incomingBlockMessages:   make(map[string]Message),
		outgoingBlockMessages:   make(map[string]Message),
		sessions:                make(map[string]Session),
	}
}

type DefaultCoapServer struct {
	messageIds            map[uint16]time.Time
	incomingBlockMessages map[string]Message
	outgoingBlockMessages map[string]Message

	routes       []*Route
	events       *Events
	observations map[string][]*Observation

	fnHandleHTTPProxy ProxyHandler
	fnHandleCOAPProxy ProxyHandler
	fnProxyFilter     ProxyFilter

	stopChannel chan int

	coapResponseChannelsMap map[uint16]chan *CoapResponseChannel

	sessions map[string]Session
}

func (s *DefaultCoapServer) GetEvents() *Events {
	return s.events
}

func (s *DefaultCoapServer) addDiscoveryRoute() {
	var discoveryRoute RouteHandler = func(req Request) Response {
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
			}
		}
		ack.Payload = payload.NewPlainTextPayload(buf.String())
		resp := message.NewResponseWithMessage(ack)
		return resp
	}
	s.NewRoute("/.well-known/core", Get, discoveryRoute)
}

func (s *DefaultCoapServer) ListenAndServeDTLS(addr string, cfg *ServerConfiguration) {
	s.addDiscoveryRoute()
	// s.serve()
}

func (s *DefaultCoapServer) ListenAndServe(addr string, cfg *ServerConfiguration) {
	s.addDiscoveryRoute()

	conn := s.createConn(addr)

	if conn == nil {
		log.Fatal("An error occured starting up CoAP Server")
	} else {
		log.Println("Started CoAP Server ", conn.LocalAddr())
		go s.handleIncomingData(conn)
		go s.events.Started(s)
		go s.handleMessageIDPurge()
	}
}

func (s *DefaultCoapServer) createConn(addr string) ServerConnection {
	// if use dTLS
	localHost := addr
	if !strings.Contains(localHost, ":") {
		localHost = ":" + localHost
	}
	localAddr, err := net.ResolveUDPAddr("udp6", localHost)
	if err != nil {
		// s.events.Error(err)
		panic(err.Error())
	}

	conn, err := net.ListenUDP(UDP, localAddr)
	if err != nil {
		// s.events.Error(err)
		panic(err.Error())
	}

	return &UDPServerConnection{
		conn: conn,
	}
}

func (s *DefaultCoapServer) handleIncomingData(conn ServerConnection) {
	readBuf := make([]byte, MaxPacketSize)
	for {
		select {
		case <-s.stopChannel:
			return

		default:
			// continue
		}

		len, addr, err := conn.ReadFrom(readBuf)
		if err == nil {
			// msgBuf := readBuf[:len]
			msgBuf := make([]byte, len)
			copy(msgBuf, readBuf[:len])

			ssn := s.sessions[addr.String()]
			if ssn == nil {
				ssn = s.createSession(addr, conn, s)
				s.sessions[addr.String()] = ssn
			}
			ssn.Write(msgBuf)
			go s.handleSession(ssn)
		} else {
			log.Println("Error occured reading UDP", err)
		}

	}
}

func (s *DefaultCoapServer) Stop() {
	// s.localConn.Close()
	close(s.stopChannel)
}

func (s *DefaultCoapServer) UpdateBlockMessageFragment(client string, msg *Message, seq uint32) {
	msgs := s.incomingBlockMessages[client]

	if msgs == nil {
		msgs = &message.BlockMessage{
			Sequence:   0,
			MessageBuf: []byte{},
		}
	}

	msgs.Sequence = seq
	msgs.MessageBuf = append(msgs.MessageBuf, msg.Payload.GetBytes()...)

	s.incomingBlockMessages[client] = msgs
}

func (s *DefaultCoapServer) FlushBlockMessagePayload(origin string) MessagePayload {
	msgs := s.incomingBlockMessages[origin]

	payload := msgs.MessageBuf

	return NewBytesPayload(payload)
}

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

func (s *DefaultCoapServer) handleSession(session Session) {
	msgBuf := session.Read()
	msg, err := BytesToMessage(msgBuf)
	if err != nil {
		panic(err.Error())
	}

	if msg.MessageType == MessageAcknowledgment {
		handleResponse(s, msg, session)
	} else {
		handleRequest(s, msg, session)
	}

	s.closeSession(session)
	// s.closeSession(session)

	// TODO: Close Session?
}

func (s *DefaultCoapServer) closeSession(ssn Session) {
	delete(s.sessions, ssn.GetAddress().String())
}

func (s *DefaultCoapServer) createSession(addr net.Addr, conn ServerConnection, server CoapServer) Session {
	return &ServerSession{
		addr:   addr,
		conn:   conn,
		server: server,
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

func (s *DefaultCoapServer) storeNewOutgoingBlockMessage(client string, payload []byte) {
	bm := NewBlockMessage()
	bm.MessageBuf = payload
	s.outgoingBlockMessages[client] = bm
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

			go SendMessage(req.GetMessage(), r.Session)
		}
	}
}

func (s *DefaultCoapServer) AddObservation(resource, token string, session Session) {
	s.observations[resource] = append(s.observations[resource], NewObservation(session, token, resource))
}

func (s *DefaultCoapServer) HasObservation(resource string, addr net.Addr) bool {
	obs := s.observations[resource]
	if obs == nil {
		return false
	}

	for _, o := range obs {
		if o.Session.GetAddress().String() == addr.String() {
			return true
		}
	}
	return false
}

func (s *DefaultCoapServer) RemoveObservation(resource string, addr net.Addr) {
	obs := s.observations[resource]
	if obs == nil {
		return
	}

	for idx, o := range obs {
		if o.Session.GetAddress().String() == addr.String() {
			s.observations[resource] = append(obs[:idx], obs[idx+1:]...)
			return
		}
	}
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

func (s *DefaultCoapServer) AllowProxyForwarding(msg *Message, addr net.Addr) bool {
	return s.fnProxyFilter(msg, addr)
}

func (s *DefaultCoapServer) ForwardCoap(msg *Message, session Session) {
	s.fnHandleCOAPProxy(s, msg, session)
}

func (s *DefaultCoapServer) ForwardHTTP(msg *Message, session Session) {
	s.fnHandleHTTPProxy(s, msg, session)
}

func (s *DefaultCoapServer) GetRoutes() []*Route {
	return s.routes
}

func (s *DefaultCoapServer) IsDuplicateMessage(msg *Message) bool {
	_, ok := s.messageIds[msg.MessageID]

	return ok
}

func (s *DefaultCoapServer) UpdateMessageTS(msg *Message) {
	s.messageIds[msg.MessageID] = time.Now()
}

func NewResponseChannel() (ch chan *CoapResponseChannel) {
	ch = make(chan *CoapResponseChannel)

	return
}

func AddResponseChannel(c CoapServer, msgId uint16, ch chan *CoapResponseChannel) {
	s := c.(*DefaultCoapServer)

	s.coapResponseChannelsMap[msgId] = ch
}

func DeleteResponseChannel(c CoapServer, msgId uint16) {
	s := c.(*DefaultCoapServer)

	delete(s.coapResponseChannelsMap, msgId)
}

func GetResponseChannel(c CoapServer, msgId uint16) (ch chan *CoapResponseChannel) {
	s := c.(*DefaultCoapServer)
	ch = s.coapResponseChannelsMap[msgId]

	return
}

func NewObservation(session Session, token string, resource string) *Observation {
	return &Observation{
		Session:     session,
		Token:       token,
		Resource:    resource,
		NotifyCount: 0,
	}
}

type Observation struct {
	Session     Session
	Token       string
	Resource    string
	NotifyCount int
}

func _doSendMessage(msg Message, session Session, ch chan *CoapResponseChannel) {
	resp := &CoapResponseChannel{}

	b, err := message.MessageToBytes(msg)
	if err != nil {
		resp.Error = err
		ch <- resp
	}

	conn := session.GetConnection()
	addr := session.GetAddress()

	_, err = conn.WriteTo(b, addr)
	session.FlushBuffer()
	if err != nil {
		resp.Error = err
		ch <- resp
	}

	if msg.MessageType == MessageNonConfirmable {
		resp.Response = NewResponse(NewEmptyMessage(msg.MessageID), nil)
		ch <- resp
	}

	AddResponseChannel(session.GetServer(), msg.MessageID, ch)
}

func SendMessage(msg *Message, session Session) (Response, error) {
	if session.GetConnection() == nil {
		return nil, ErrNilConn
	}

	if msg == nil {
		return nil, ErrNilMessage
	}

	if session.GetAddress() == nil {
		return nil, ErrNilAddr
	}

	ch := server.NewResponseChannel()
	go _doSendMessage(msg, session, ch)
	respCh := <-ch

	return respCh.Response, respCh.Error
}

type CoapResponseChannel struct {
	Response Response
	Error    error
}
