package goap

import (
	"bytes"
	"log"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func NewCoapServer(localAddr *net.UDPAddr, remoteAddr *net.UDPAddr) *CoapServer {

	return &CoapServer{
		remoteAddr: remoteAddr,
		localAddr:  localAddr,
	}
}

type CoapServer struct {
	localAddr  *net.UDPAddr
	remoteAddr *net.UDPAddr
	conn       *net.UDPConn
	messageIds map[uint16]time.Time
	routes     []*Route

	evtOnStartup  EventHandler
	evtOnClose    EventHandler
	evtOnDiscover EventHandler
	evtOnError    EventHandler
	evtOnMessage  EventHandler
}

func (s *CoapServer) Start() {

	var discoveryRoute RouteHandler = func(req *CoapRequest) *CoapResponse {
		msg := req.GetMessage()

		ack := NewMessageOfType(TYPE_ACKNOWLEDGEMENT, msg.MessageId)
		ack.Code = COAPCODE_205_CONTENT
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

		/*
		   if s.fnEventDiscover != nil {
		       e := NewEvent()
		       e.Message = ack

		       ack = s.fnEventDiscover(e)
		   }
		*/

		resp := NewResponseWithMessage(ack)

		return resp
	}

	s.NewRoute(".well-known/core", GET, discoveryRoute)
	s.serveServer()
}

func (s *CoapServer) serveServer() {
	s.messageIds = make(map[uint16]time.Time)

	conn, err := net.ListenUDP("udp", s.localAddr)
	IfErr(err)
	s.conn = conn

	log.Println("Started server ", conn.LocalAddr())

	CallEvent(s.evtOnStartup, EmptyEventPayload())

	s.handleMessageIdPurge()

	readBuf := make([]byte, BUF_SIZE)
	for {
		len, addr, err := conn.ReadFromUDP(readBuf)

		if err == nil {

			msgBuf := make([]byte, len)
			copy(msgBuf, readBuf)

			// Look for route handler matching path and then dispatch
			go s.handleMessage(msgBuf, conn, addr)
		}
	}
}

func (s *CoapServer) handleMessageIdPurge() {
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

func (s *CoapServer) handleMessage(msgBuf []byte, conn *net.UDPConn, addr *net.UDPAddr) {
	msg, err := BytesToMessage(msgBuf)

	CallEvent(s.evtOnMessage, EmptyEventPayload())

	// Unsupported Method
	if msg.Code != GET && msg.Code != POST && msg.Code != PUT && msg.Code != DELETE {
		ret := NewMessage(TYPE_NONCONFIRMABLE, COAPCODE_501_NOT_IMPLEMENTED, msg.MessageId)
		ret.CloneOptions(msg, OPTION_URI_PATH, OPTION_CONTENT_FORMAT)

		SendMessageTo(ret, conn, addr)
		return
	}

	if err != nil {
		if err == ERR_UNKNOWN_CRITICAL_OPTION {
			if msg.MessageType == TYPE_CONFIRMABLE {
				SendError402BadOption(msg.MessageId, conn, addr)
				return
			} else {
				// Ignore silently
				return
			}
		}
	}

	route, attrs, err := MatchingRoute(msg, s.routes)
	if err != nil {
		if err == ERR_NO_MATCHING_ROUTE {
			ret := NewMessage(TYPE_NONCONFIRMABLE, COAPCODE_404_NOT_FOUND, msg.MessageId)
			ret.CloneOptions(msg, OPTION_URI_PATH, OPTION_CONTENT_FORMAT)
			ret.Token = msg.Token

			SendMessageTo(ret, conn, addr)
			CallEvent(s.evtOnError, EmptyEventPayload())
			return
		}

		if err == ERR_NO_MATCHING_METHOD {
			ret := NewMessage(TYPE_NONCONFIRMABLE, COAPCODE_405_METHOD_NOT_ALLOWED, msg.MessageId)
			ret.CloneOptions(msg, OPTION_URI_PATH, OPTION_CONTENT_FORMAT)

			SendMessageTo(ret, conn, addr)
			CallEvent(s.evtOnError, EmptyEventPayload())
			return
		}

		if err == ERR_UNSUPPORTED_CONTENT_FORMAT {
			ret := NewMessage(TYPE_NONCONFIRMABLE, COAPCODE_415_UNSUPPORTED_CONTENT_FORMAT, msg.MessageId)
			ret.CloneOptions(msg, OPTION_URI_PATH, OPTION_CONTENT_FORMAT)

			SendMessageTo(ret, conn, addr)
			CallEvent(s.evtOnError, EmptyEventPayload())
			return
		}
	}

	// Duplicate Message ID Check
	_, dupe := s.messageIds[msg.MessageId]
	if dupe {
		log.Println("Duplicate Message ID ", msg.MessageId)
		if msg.MessageType == TYPE_CONFIRMABLE {
			ret := NewMessage(TYPE_RESET, COAPCODE_0_EMPTY, msg.MessageId)
			ret.CloneOptions(msg, OPTION_URI_PATH, OPTION_CONTENT_FORMAT)

			SendMessageTo(ret, conn, addr)
		}
		return
	}

	if err == nil {
		s.messageIds[msg.MessageId] = time.Now()

		// TODO: #47 - Forward Proxy

		// Auto acknowledge
		if msg.MessageType == TYPE_CONFIRMABLE && route.AutoAck {
			ack := NewMessageOfType(TYPE_ACKNOWLEDGEMENT, msg.MessageId)

			SendMessageTo(ack, conn, addr)
		}

		req := NewRequestFromMessage(msg, attrs)

		resp := route.Handler(req)

		// TODO: Validate Message before sending (.e.g missing messageId)
		SendMessageTo(resp.GetMessage(), conn, addr)
	}
}

func (s *CoapServer) NewRoute(path string, method CoapCode, fn RouteHandler) *Route {
	re, _ := regexp.Compile(`{[a-z]+}`)
	matches := re.FindAllStringSubmatch(path, -1)

	path = "^" + path
	for _, b := range matches {
		origAttr := b[0]
		attr := strings.Replace(strings.Replace(origAttr, "{", "", -1), "}", "", -1)
		frag := `(?P<` + attr + `>\w+)`
		path = strings.Replace(path, origAttr, frag, -1)
	}
	path += "$"
	re, _ = regexp.Compile(path)

	/*
	   OnNewRoute
	       Get all values between #{ }
	       Construct New RegEx
	           Create SubGroups
	           Escape any RegEx Values
	       Compile and Store Compiled RegEx

	*/

	r := &Route{
		AutoAck: false,
		Path:    path,
		Method:  method,
		Handler: fn,
		RegEx:   re,
	}
	s.routes = append(s.routes, r)

	return r
}

func (c *CoapServer) Send(req *CoapRequest) (*CoapResponse, error) {
	return SendMessageTo(req.GetMessage(), c.conn, c.remoteAddr)
}

func (c *CoapServer) SendTo(req *CoapRequest, addr *net.UDPAddr) (*CoapResponse, error) {
	return SendMessageTo(req.GetMessage(), c.conn, addr)
}

func (c *CoapServer) OnStartup(eh EventHandler) {
	c.evtOnStartup = eh
}

func (c *CoapServer) OnClose(eh EventHandler) {
	c.evtOnStartup = eh
}

func (c *CoapServer) OnDiscover(eh EventHandler) {
	c.evtOnStartup = eh
}

func (c *CoapServer) OnError(eh EventHandler) {
	c.evtOnStartup = eh
}

func (c *CoapServer) OnMessage(eh EventHandler) {
	c.evtOnStartup = eh
}
