package canopus

import (
	"bytes"
	"log"
	"net"
	"strconv"
	"time"
)

type ProxyType int

const (
	PROXY_HTTP ProxyType = 0
	PROXY_COAP ProxyType = 1
)

func NewLocalServer() CoapServer {
	localAddr, err := net.ResolveUDPAddr("udp6", ":5683")
	if err != nil {
		log.Fatal("Error starting CoAP Server: ", err)
		return nil
	}
	return NewServer(localAddr, nil)
}

func NewCoapServer(local string) CoapServer {
	localAddr, err := net.ResolveUDPAddr("udp6", local)

	if err != nil {
		log.Println("Error instantiating CoAP Server")
		return nil
	}

	return NewServer(localAddr, nil)
}

func NewCoapClient() CoapServer {
	localAddr, _ := net.ResolveUDPAddr("udp6", ":0")

	return NewServer(localAddr, nil)
}

func NewServer(localAddr *net.UDPAddr, remoteAddr *net.UDPAddr) CoapServer {
	return &DefaultCoapServer{
		remoteAddr:            remoteAddr,
		localAddr:             localAddr,
		events:                NewCanopusEvents(),
		observations:          make(map[string][]*Observation),
		fnHandleCoapProxy: 		NullProxyHandler,
		fnHandleHttpProxy: 		NullProxyHandler,
		fnProxyFilter:			NullProxyFilter,
		stopChannel:           make(chan int),
	}
}

type DefaultCoapServer struct {
	localAddr  *net.UDPAddr
	remoteAddr *net.UDPAddr

	localConn  *net.UDPConn
	remoteConn *net.UDPConn

	messageIds   map[uint16]time.Time
	routes       []*Route
	events       *CanopusEvents
	observations map[string][]*Observation

	fnHandleHttpProxy ProxyHandler
	fnHandleCoapProxy ProxyHandler
	fnProxyFilter	  ProxyFilter

	stopChannel chan int
}

func (s *DefaultCoapServer) GetEvents() *CanopusEvents {
	return s.events
}

func (s *DefaultCoapServer) Start() {

	var discoveryRoute RouteHandler = func(req CoapRequest) CoapResponse {
		msg := req.GetMessage()

		ack := ContentMessage(msg.MessageId, TYPE_ACKNOWLEDGEMENT)
		ack.Token = make([]byte, len(msg.Token))
		copy(ack.Token, msg.Token)

		ack.AddOption(OPTION_CONTENT_FORMAT, MEDIATYPE_APPLICATION_LINK_FORMAT)

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

	s.NewRoute("/.well-known/core", GET, discoveryRoute)
	s.serveServer()
}

func (s *DefaultCoapServer) serveServer() {
	s.messageIds = make(map[uint16]time.Time)

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
	s.handleMessageIdPurge()

	readBuf := make([]byte, MAX_PACKET_SIZE)
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

func (s *DefaultCoapServer) handleMessageIdPurge() {
	// Routine for clearing up message IDs which has expired
	ticker := time.NewTicker(MESSAGEID_PURGE_DURATION * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				for k, v := range s.messageIds {
					elapsed := time.Since(v)
					if elapsed > MESSAGEID_PURGE_DURATION {
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

	if msg.MessageType == TYPE_ACKNOWLEDGEMENT {
		handleResponse(s, msg, conn, addr)
	} else {
		handleRequest(s, err, msg, conn,addr)
	}
}

func (s *DefaultCoapServer) Get(path string, fn RouteHandler) *Route {
	return s.add(METHOD_GET, path, fn)
}

func (s *DefaultCoapServer) Delete(path string, fn RouteHandler) *Route {
	return s.add(METHOD_DELETE, path, fn)
}

func (s *DefaultCoapServer) Put(path string, fn RouteHandler) *Route {
	return s.add(METHOD_PUT, path, fn)
}

func (s *DefaultCoapServer) Post(path string, fn RouteHandler) *Route {
	return s.add(METHOD_POST, path, fn)
}

func (s *DefaultCoapServer) Options(path string, fn RouteHandler) *Route {
	return s.add(METHOD_OPTIONS, path, fn)
}

func (s *DefaultCoapServer) Patch(path string, fn RouteHandler) *Route {
	return s.add(METHOD_PATCH, path, fn)
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

func (c *DefaultCoapServer) Send(req CoapRequest) (CoapResponse, error) {
	c.events.Message(req.GetMessage(), false)
	response, err := SendMessageTo(req.GetMessage(), NewCanopusUDPConnection(c.localConn), c.remoteAddr)

	if err != nil {
		c.events.Error(err)
		return response, err
	}
	c.events.Message(response.GetMessage(), true)

	return response, err
}

func (c *DefaultCoapServer) SendTo(req CoapRequest, addr *net.UDPAddr) (CoapResponse, error) {
	return SendMessageTo(req.GetMessage(), NewCanopusUDPConnection(c.localConn), addr)
}

func (c *DefaultCoapServer) NotifyChange(resource, value string, confirm bool) {
	t := c.observations[resource]

	if t != nil {
		var req CoapRequest

		if confirm {
			req = NewRequest(TYPE_CONFIRMABLE, COAPCODE_205_CONTENT, GenerateMessageId())
		} else {
			req = NewRequest(TYPE_ACKNOWLEDGEMENT, COAPCODE_205_CONTENT, GenerateMessageId())
		}

		for _, r := range t {
			req.SetToken(r.Token)
			req.SetStringPayload(value)
			req.SetRequestURI(r.Resource)
			r.NotifyCount++
			req.GetMessage().AddOption(OPTION_OBSERVE, r.NotifyCount)

			go c.SendTo(req, r.Addr)
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

func (c *DefaultCoapServer) Dial(host string) {
	c.Dial6(host)
}

func (c *DefaultCoapServer) Dial6(host string) {
	remoteAddr, _ := net.ResolveUDPAddr("udp6", host)

	c.remoteAddr = remoteAddr
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

func (s *DefaultCoapServer) ProxyHttp(enabled bool) {
	if enabled {
		s.fnHandleHttpProxy = HttpProxyHandler
	} else {
		s.fnHandleHttpProxy = NullProxyHandler
	}
}

func (s *DefaultCoapServer) ProxyCoap(enabled bool) {
	if enabled {
		s.fnHandleCoapProxy = CoapProxyHandler
	} else {
		s.fnHandleCoapProxy = NullProxyHandler
	}
}

func (s *DefaultCoapServer) AllowProxyForwarding(msg *Message, addr *net.UDPAddr) (bool) {
	return s.fnProxyFilter(msg, addr)
}

func (s *DefaultCoapServer) ForwardCoap(msg *Message, conn *net.UDPConn, addr *net.UDPAddr) {
	s.fnHandleCoapProxy(msg, conn, addr)
}

func (s *DefaultCoapServer) ForwardHttp(msg *Message, conn *net.UDPConn, addr *net.UDPAddr) {
	s.fnHandleHttpProxy(msg, conn, addr)
}

func (s *DefaultCoapServer) GetRoutes() []*Route {
	return s.routes
}

func (s *DefaultCoapServer) GetLocalAddress() *net.UDPAddr {
	return s.localAddr
}

func (s *DefaultCoapServer) IsDuplicateMessage(msg *Message) bool {
	_, ok := s.messageIds[msg.MessageId]

	return ok
}

func (s *DefaultCoapServer) UpdateMessageTS(msg *Message) {
	s.messageIds[msg.MessageId] = time.Now()
}

////////////////////////////////////////////////////////////////////////////////
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
