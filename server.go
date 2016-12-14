package canopus

import (
	"bytes"
	"crypto/rand"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

var DTLS_SERVER_SESSIONS = make(map[int32]*DTLSServerSession)
var NEXT_SESSION_ID int32 = 0
var DTLS_CLIENT_CONNECTIONS = make(map[int32]*DTLSConnection)

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
		createdSession:          make(chan Session),
	}
}

type DefaultCoapServer struct {
	messageIds            map[uint16]time.Time
	incomingBlockMessages map[string]Message
	outgoingBlockMessages map[string]Message

	routes       []Route
	events       Events
	observations map[string][]*Observation

	fnHandleHTTPProxy ProxyHandler
	fnHandleCOAPProxy ProxyHandler
	fnProxyFilter     ProxyFilter

	stopChannel chan int

	coapResponseChannelsMap map[uint16]chan *CoapResponseChannel

	sessions       map[string]Session
	createdSession chan Session
	serverConfig   *ServerConfiguration

	cookieSecret []byte

	fnPskHandler func(id string) []byte
}

func (s *DefaultCoapServer) DeleteSession(ssn Session) {
	s.closeSession(ssn)
}

func (s *DefaultCoapServer) HandlePSK(fn func(id string) []byte) {
	s.fnPskHandler = fn
}

func (s *DefaultCoapServer) handleRequest(msg Message, session Session) {
	if msg.GetMessageType() != MessageReset {
		// Unsupported Method
		if msg.GetCode() != Get && msg.GetCode() != Post && msg.GetCode() != Put && msg.GetCode() != Delete {
			s.handleReqUnsupportedMethodRequest(msg, session)
			return
		}

		// Proxy
		if IsProxyRequest(msg) {
			s.handleReqProxyRequest(msg, session)
		} else {
			route, attrs, err := MatchingRoute(msg.GetURIPath(), MethodString(msg.GetCode()), msg.GetOptions(OptionContentFormat), s.GetRoutes())
			if err != nil {
				s.GetEvents().Error(err)
				if err == ErrNoMatchingRoute {
					s.handleReqNoMatchingRoute(msg, session)
					return
				}

				if err == ErrNoMatchingMethod {
					s.handleReqNoMatchingMethod(msg, session)
					return
				}

				if err == ErrUnsupportedContentFormat {
					s.handleReqUnsupportedContentFormat(msg, session)
					return
				}

				log.Println("Error occured parsing inbound message")
				return
			}

			// Duplicate Message ID Check
			if s.isDuplicateMessage(msg) {
				PrintMessage(msg)
				if msg.GetMessageType() == MessageConfirmable {
					log.Println("Duplicate Message ID ", msg.GetMessageId())
					s.handleReqDuplicateMessageID(msg, session)
				}
				return
			}

			s.updateMessageTS(msg)

			// Auto acknowledge
			// TODO: Necessary?
			if msg.GetMessageType() == MessageConfirmable && route.AutoAcknowledge() {
				s.handleRequestAcknowledge(msg, session)
			}
			req := NewClientRequestFromMessage(msg, attrs, session)

			if msg.GetMessageType() == MessageConfirmable {
				// Observation Request
				obsOpt := msg.GetOption(OptionObserve)
				if obsOpt != nil {
					s.handleReqObserve(req, msg, session)
				}
			}
			opt := req.GetMessage().GetOption(OptionBlock1)
			if opt != nil {
				blockOpt := Block1OptionFromOption(opt)

				// 0000 1 010
				/*
									[NUM][M][SZX]
									2 ^ (2 + 4)
									2 ^ 6 = 32
									Size = 2 ^ (SZX + 4)

									The value 7 for SZX (which would
					      	indicate a block size of 2048) is reserved, i.e. MUST NOT be sent
					      	and MUST lead to a 4.00 Bad Request response code upon reception
					      	in a request.
				*/

				if blockOpt.Value != nil {
					if blockOpt.Code == OptionBlock1 {
						exp := blockOpt.Exponent()

						if exp == 7 {
							s.handleReqBadRequest(msg, session)
							return
						}

						// szx := blockOpt.Size()
						hasMore := blockOpt.HasMore()
						seqNum := blockOpt.Sequence()
						// fmt.Println("Out Values == ", blockOpt.Value, exp, szx, 2, hasMore, seqNum)

						s.GetEvents().BlockMessage(msg, true)

						s.updateBlockMessageFragment(session.GetAddress().String(), msg, seqNum)

						if hasMore {
							s.handleReqContinue(msg, session)

							// Auto Respond to client

						} else {
							// TODO: Check if message is too large
							msg = NewMessage(msg.GetMessageType(), msg.GetCode(), msg.GetMessageId())
							msg.SetPayload(s.flushBlockMessagePayload(session.GetAddress().String()))
							req = NewClientRequestFromMessage(msg, attrs, session)
						}
					} else if blockOpt.Code == OptionBlock2 {

					} else {
						// TOOO: Invalid Block option Code
					}
				}
			}

			resp := route.Handle(req)
			_, nilresponse := resp.(NilResponse)
			if !nilresponse {
				respMsg := resp.GetMessage().(*CoapMessage)
				respMsg.SetToken(req.GetMessage().GetToken())

				// TODO: Validate Message before sending (e.g missing messageId)
				err := ValidateMessage(respMsg)
				if err == nil {
					s.GetEvents().Message(respMsg, false)
					SendMessage(respMsg, session)
				}
			}
		}
	}
}

func (s *DefaultCoapServer) handleReqObserve(req Request, msg Message, session Session) {
	// TODO: if server doesn't allow observing, return error
	addr := session.GetAddress()

	// TODO: Check if observation has been registered, if yes, remove it (observation == cancel)
	resource := msg.GetURIPath()

	if s.HasObservation(resource, addr) {
		// Remove observation of client
		s.RemoveObservation(resource, addr)

		// Observe Cancel Request & Fire OnObserveCancel Event
		s.GetEvents().ObserveCancelled(resource, msg)
	} else {
		// Register observation of client
		s.AddObservation(msg.GetURIPath(), string(msg.GetToken()), session)

		// Observe Request & Fire OnObserve Event
		s.GetEvents().Observe(resource, msg)
	}

	req.GetMessage().AddOption(OptionObserve, 1)
}

func (s *DefaultCoapServer) handleResponse(msg Message, session Session) {
	defer s.closeSession(session)
	if msg.GetOption(OptionObserve) != nil {
		s.handleAcknowledgeObserveRequest(msg)
		return
	}

	ch := GetResponseChannel(s, msg.GetMessageId())
	if ch != nil {
		resp := &CoapResponseChannel{
			Response: NewResponse(msg, nil),
		}
		ch <- resp
		DeleteResponseChannel(s, msg.GetMessageId())
	}
}

func (s *DefaultCoapServer) GetEvents() Events {
	return s.events
}

func (s *DefaultCoapServer) addDiscoveryRoute() {
	var discoveryRoute RouteHandler = func(req Request) Response {
		msg := req.GetMessage()

		var buf bytes.Buffer
		for _, r := range s.routes {
			if r.GetConfiguredPath() != ".well-known/core" {
				buf.WriteString("</")
				buf.WriteString(r.GetConfiguredPath())
				buf.WriteString(">")

				// Media Types
				lenMt := len(r.GetMediaTypes())
				if lenMt > 0 {
					buf.WriteString(";ct=")
					for idx, mt := range r.GetMediaTypes() {

						buf.WriteString(strconv.Itoa(int(mt)))
						if idx+1 < lenMt {
							buf.WriteString(" ")
						}
					}
				}

				buf.WriteString(",")
			}
		}

		ack := ContentMessage(msg.GetMessageId(), MessageAcknowledgment)
		ack.SetToken(msg.GetToken())
		ack.SetPayload(NewPlainTextPayload(buf.String()))
		ack.AddOption(OptionContentFormat, MediaTypeApplicationLinkFormat)
		resp := NewResponseWithMessage(ack)

		return resp
	}
	s.NewRoute("/.well-known/core", Get, discoveryRoute)
}

func (s *DefaultCoapServer) ListenAndServeDTLS(addr string) {
	s.addDiscoveryRoute()

	conn := s.createConn(addr)

	ctx, err := NewServerDtlsContext()
	if err != nil {
		panic("Unable to create SSL Context:" + err.Error())
	}

	if conn == nil {
		log.Fatal("An error occured starting up CoAPS Server")
	} else {
		secret := make([]byte, 32)
		if n, err := rand.Read(secret); n != 32 || err != nil {
			panic(err)
		}

		s.cookieSecret = secret
		log.Println("Started CoAPS Server ", conn.LocalAddr())
		go s.handleIncomingDTLSData(conn, ctx)
		go s.events.Started(s)
		go s.handleMessageIDPurge()
	}
}

func (s *DefaultCoapServer) ListenAndServe(addr string) {
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
	localHost := addr
	if !strings.Contains(localHost, ":") {
		localHost = ":" + localHost
	}
	localAddr, err := net.ResolveUDPAddr("udp6", localHost)
	if err != nil {
		panic(err.Error())
	}

	conn, err := net.ListenUDP(UDP, localAddr)
	if err != nil {
		panic(err.Error())
	}

	return &UDPServerConnection{
		conn: conn,
	}
}

func (s *DefaultCoapServer) handleIncomingDTLSData(conn ServerConnection, ctx *ServerDtlsContext) {
	readBuf := make([]byte, MaxPacketSize)
	go func() {
		for {
			select {
			case <-s.stopChannel:
				return

			default:
				// continue
			}

			len, addr, err := conn.ReadFrom(readBuf)
			if err == nil {
				msgBuf := make([]byte, len)
				copy(msgBuf, readBuf[:len])
				ssn := s.sessions[addr.String()]
				if ssn == nil {
					ssn = &DTLSServerSession{
						UDPServerSession: UDPServerSession{
							addr:   addr,
							conn:   conn,
							server: s,
							buf:    []byte{},
							rcvd:   make(chan []byte, 1),
						},
					}
					err := newSslSession(ssn.(*DTLSServerSession), ctx, s.fnPskHandler)
					if err != nil {
						panic(err.Error())
					}
					s.sessions[addr.String()] = ssn
					s.createdSession <- ssn
				}

				ssn.(*DTLSServerSession).rcvd <- msgBuf
			} else {
				logMsg("Error occured reading UDP", err)
			}
		}
	}()

	go func() {
		for {
			ssn := <-s.createdSession
			go s.handleSession(ssn)
		}
	}()

}

func (s *DefaultCoapServer) handleIncomingData(conn ServerConnection) {
	readBuf := make([]byte, MaxPacketSize)
	go func() {
		for {
			select {
			case <-s.stopChannel:
				return

			default:
				// continue
			}

			len, addr, err := conn.ReadFrom(readBuf)
			if err == nil {
				msgBuf := make([]byte, len)
				copy(msgBuf, readBuf[:len])
				ssn := s.sessions[addr.String()]
				if ssn == nil {
					ssn = &UDPServerSession{
						addr:   addr,
						conn:   conn,
						server: s,
						rcvd:   make(chan []byte),
					}
					if err != nil {
						panic(err.Error())
					}
					s.sessions[addr.String()] = ssn
				}
				go func() {
					ssn.(*UDPServerSession).rcvd <- msgBuf
				}()
				go s.handleSession(ssn)
			} else {
				logMsg("Error occured reading UDP", err)
			}
		}
	}()

}

func (s *DefaultCoapServer) GetSession(addr string) Session {
	return s.sessions[addr]
}

func (s *DefaultCoapServer) Stop() {
	close(s.stopChannel)
}

func (s *DefaultCoapServer) updateBlockMessageFragment(client string, msg Message, seq uint32) {
	msgs := s.incomingBlockMessages[client]

	if msgs == nil {
		msgs = &CoapBlockMessage{
			Sequence:   0,
			MessageBuf: []byte{},
		}
	}

	blockMsgs := msgs.(*CoapBlockMessage)
	blockMsgs.Sequence = seq
	blockMsgs.MessageBuf = append(blockMsgs.MessageBuf, msg.GetPayload().GetBytes()...)

	s.incomingBlockMessages[client] = msgs
}

func (s *DefaultCoapServer) flushBlockMessagePayload(origin string) MessagePayload {
	msgs := s.incomingBlockMessages[origin]

	blockMsg := msgs.(*CoapBlockMessage)
	payload := blockMsg.MessageBuf

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

func (s *DefaultCoapServer) GetCookieSecret() []byte {
	return s.cookieSecret
}

func (s *DefaultCoapServer) handleSession(session Session) {
	msgBuf := make([]byte, 1500)
	n, _ := session.Read(msgBuf)

	msg, err := BytesToMessage(msgBuf[:n])
	if err != nil {
		logMsg(err.Error())
		s.handleReqBadRequest(msg, session)
	}

	if msg.GetMessageType() == MessageAcknowledgment {
		s.handleResponse(msg, session)
	} else {
		s.handleRequest(msg, session)
	}
}

func (s *DefaultCoapServer) closeSession(ssn Session) {
	delete(s.sessions, ssn.GetAddress().String())
}

func (s *DefaultCoapServer) Get(path string, fn RouteHandler) Route {
	return s.add(MethodGet, path, fn)
}

func (s *DefaultCoapServer) Delete(path string, fn RouteHandler) Route {
	return s.add(MethodDelete, path, fn)
}

func (s *DefaultCoapServer) Put(path string, fn RouteHandler) Route {
	return s.add(MethodPut, path, fn)
}

func (s *DefaultCoapServer) Post(path string, fn RouteHandler) Route {
	return s.add(MethodPost, path, fn)
}

func (s *DefaultCoapServer) Options(path string, fn RouteHandler) Route {
	return s.add(MethodOptions, path, fn)
}

func (s *DefaultCoapServer) Patch(path string, fn RouteHandler) Route {
	return s.add(MethodPatch, path, fn)
}

func (s *DefaultCoapServer) add(method string, path string, fn RouteHandler) Route {
	route := CreateNewRegExRoute(path, method, fn)
	s.routes = append(s.routes, route)

	return route
}

func (s *DefaultCoapServer) NewRoute(path string, method CoapCode, fn RouteHandler) Route {
	route := CreateNewRegExRoute(path, MethodString(method), fn)
	s.routes = append(s.routes, route)

	return route
}

func (s *DefaultCoapServer) storeNewOutgoingBlockMessage(client string, payload []byte) {
	bm := NewBlockMessage().(*CoapBlockMessage)
	bm.MessageBuf = payload
	s.outgoingBlockMessages[client] = bm
}

func (s *DefaultCoapServer) NotifyChange(resource, value string, confirm bool) {
	t := s.observations[resource]

	if t != nil {
		var req Request

		if confirm {
			req = NewRequest(MessageConfirmable, CoapCodeContent)
		} else {
			req = NewRequest(MessageAcknowledgment, CoapCodeContent)
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

func (s *DefaultCoapServer) ProxyOverHttp(enabled bool) {
	if enabled {
		s.fnHandleHTTPProxy = HTTPProxyHandler
	} else {
		s.fnHandleHTTPProxy = NullProxyHandler
	}
}

func (s *DefaultCoapServer) ProxyOverCoap(enabled bool) {
	if enabled {
		s.fnHandleCOAPProxy = COAPProxyHandler
	} else {
		s.fnHandleCOAPProxy = NullProxyHandler
	}
}

func (s *DefaultCoapServer) AllowProxyForwarding(msg Message, addr net.Addr) bool {
	return s.fnProxyFilter(msg, addr)
}

func (s *DefaultCoapServer) ForwardCoap(msg Message, session Session) {
	s.fnHandleCOAPProxy(s, msg, session)
}

func (s *DefaultCoapServer) ForwardHTTP(msg Message, session Session) {
	s.fnHandleHTTPProxy(s, msg, session)
}

func (s *DefaultCoapServer) GetRoutes() []Route {
	return s.routes
}

func (s *DefaultCoapServer) isDuplicateMessage(msg Message) bool {
	_, ok := s.messageIds[msg.GetMessageId()]

	return ok
}

func (s *DefaultCoapServer) updateMessageTS(msg Message) {
	s.messageIds[msg.GetMessageId()] = time.Now()
}

func (s *DefaultCoapServer) handleReqUnknownCriticalOption(msg Message, session Session) {
	if msg.GetMessageType() == MessageConfirmable {
		SendMessage(BadOptionMessage(msg.GetMessageId(), MessageAcknowledgment), session)
	}
	return
}

func (s *DefaultCoapServer) handleReqBadRequest(msg Message, session Session) {
	if msg.GetMessageType() == MessageConfirmable {
		SendMessage(BadRequestMessage(msg.GetMessageId(), msg.GetMessageType()), session)
	}
	return
}

func (s *DefaultCoapServer) handleReqContinue(msg Message, session Session) {
	if msg.GetMessageType() == MessageConfirmable {
		SendMessage(ContinueMessage(msg.GetMessageId(), msg.GetMessageType()), session)
	}
	return
}

func (s *DefaultCoapServer) handleReqUnsupportedMethodRequest(msg Message, session Session) {
	ret := NotImplementedMessage(msg.GetMessageId(), MessageAcknowledgment)
	ret.CloneOptions(msg, OptionURIPath, OptionContentFormat)

	// c.GetEvents().Message(ret, false)
	SendMessage(ret, session)
}

func (s *DefaultCoapServer) handleReqProxyRequest(msg Message, session Session) {
	if !s.AllowProxyForwarding(msg, session.GetAddress()) {
		SendMessage(ForbiddenMessage(msg.GetMessageId(), MessageAcknowledgment), session)
	}

	proxyURI := msg.GetOption(OptionProxyURI).StringValue()
	if IsCoapURI(proxyURI) {
		s.ForwardCoap(msg, session)
	} else if IsHTTPURI(proxyURI) {
		s.ForwardHTTP(msg, session)
	} else {
		//
	}
}

func (s *DefaultCoapServer) handleReqNoMatchingRoute(msg Message, session Session) {
	ret := NotFoundMessage(msg.GetMessageId(), MessageAcknowledgment, msg.GetToken())
	ret.CloneOptions(msg, OptionURIPath, OptionContentFormat)

	SendMessage(ret, session)
}

func (s *DefaultCoapServer) handleReqNoMatchingMethod(msg Message, session Session) {
	ret := MethodNotAllowedMessage(msg.GetMessageId(), MessageAcknowledgment)
	ret.CloneOptions(msg, OptionURIPath, OptionContentFormat)

	SendMessage(ret, session)
}

func (s *DefaultCoapServer) handleReqUnsupportedContentFormat(msg Message, session Session) {
	ret := UnsupportedContentFormatMessage(msg.GetMessageId(), MessageAcknowledgment)
	ret.CloneOptions(msg, OptionURIPath, OptionContentFormat)

	// s.GetEvents().Message(ret, false)
	SendMessage(ret, session)
}

func (s *DefaultCoapServer) handleReqDuplicateMessageID(msg Message, session Session) {
	ret := EmptyMessage(msg.GetMessageId(), MessageReset)
	ret.CloneOptions(msg, OptionURIPath, OptionContentFormat)

	SendMessage(ret, session)
}

func (s *DefaultCoapServer) handleRequestAcknowledge(msg Message, session Session) {
	ack := NewMessageOfType(MessageAcknowledgment, msg.GetMessageId(), nil)

	SendMessage(ack, session)
}

func (s *DefaultCoapServer) handleAcknowledgeObserveRequest(msg Message) {
	s.GetEvents().Notify(msg.GetURIPath(), msg.GetPayload(), msg)
}

func (s *DefaultCoapServer) handleAcknowledgeObserveRequestGetSession(addr string) Session {
	return s.sessions[addr]
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

	b, err := MessageToBytes(msg)
	if err != nil {
		resp.Error = err
		ch <- resp
	}

	_, err = session.Write(b)
	if err != nil {
		resp.Error = err
		ch <- resp
	}

	if msg.GetMessageType() == MessageNonConfirmable {
		resp.Response = NewResponse(NewEmptyMessage(msg.GetMessageId()), nil)
		ch <- resp
	}
	AddResponseChannel(session.GetServer(), msg.GetMessageId(), ch)
}

func SendMessage(msg Message, session Session) (Response, error) {
	if session.GetConnection() == nil {
		return nil, ErrNilConn
	}

	if msg == nil {
		return nil, ErrNilMessage
	}

	if session.GetAddress() == nil {
		return nil, ErrNilAddr
	}

	ch := NewResponseChannel()
	go _doSendMessage(msg, session, ch)

	session.GetServer().DeleteSession(session)
	respCh := <-ch
	return respCh.Response, respCh.Error
}

type CoapResponseChannel struct {
	Response Response
	Error    error
}
