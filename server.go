package canopus

import (
	"bytes"
	"github.com/zubairhamed/go-commons/logging"
	. "github.com/zubairhamed/go-commons/network"
	"log"
	"net"
	"strconv"
	"time"
)

func NewLocalServer() *CoapServer {
	localAddr, err := net.ResolveUDPAddr("udp", ":5683")
	if err != nil {
		logging.LogError("Error starting CoAP Server: ", err)
	}
	return NewServer(localAddr, nil)
}

func NewCoapServer(local string) *CoapServer {
	localAddr, _ := net.ResolveUDPAddr("udp", local)

	return NewServer(localAddr, nil)
}

func NewServer(localAddr *net.UDPAddr, remoteAddr *net.UDPAddr) *CoapServer {
	return &CoapServer{
		remoteAddr: remoteAddr,
		localAddr:  localAddr,
		events:     make(map[EventCode]FnCanopusEvent),
		observations: make(map[string][]*Observation),
	}
}

type CoapServer struct {
	localAddr  		*net.UDPAddr
	localConn       *net.UDPConn

	remoteAddr 		*net.UDPAddr
	remoteConn		*net.UDPConn
	messageIds 		map[uint16]time.Time
	routes     		[]*Route
	events     		map[EventCode]FnCanopusEvent
	observations	map[string][]*Observation
}

func (s *CoapServer) Start() {

	var discoveryRoute RouteHandler = func(req Request) Response {
		msg := req.(*CoapRequest).GetMessage()

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
	logging.LogError(err)

	s.localConn = conn

	if conn == nil {
		log.Fatal("An error occured starting up CoAP Server")
	} else {
		log.Println("Started CoAP Server ", conn.LocalAddr())
	}

	CallEvent(EVT_START, s.events[EVT_START])

	s.handleMessageIdPurge()

	readBuf := make([]byte, BUF_SIZE)
	for {
		len, addr, err := conn.ReadFromUDP(readBuf)

		if err == nil {

			msgBuf := make([]byte, len)
			copy(msgBuf, readBuf)


			log.Println("Handling MSG Buf", msgBuf)
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
	PrintMessage(msg)

	CallEvent(EVT_MESSAGE, s.events[EVT_MESSAGE])

	if msg.MessageType != TYPE_ACKNOWLEDGEMENT && msg.MessageType != TYPE_RESET {

		// Unsupported Method
		if msg.Code != GET && msg.Code != POST && msg.Code != PUT && msg.Code != DELETE {
			log.Println("Unsupported Method ", msg.Code)
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

		route, attrs, err := MatchingRoute(msg.GetUriPath(), MethodString(msg.Code), msg.GetOptions(OPTION_CONTENT_FORMAT), s.routes)
		if err != nil {
			if err == ERR_NO_MATCHING_ROUTE {
				ret := NewMessage(TYPE_NONCONFIRMABLE, COAPCODE_404_NOT_FOUND, msg.MessageId)
				ret.CloneOptions(msg, OPTION_URI_PATH, OPTION_CONTENT_FORMAT)
				ret.Token = msg.Token

				SendMessageTo(ret, conn, addr)
				CallEvent(EVT_ERROR, s.events[EVT_ERROR])
				return
			}

			if err == ERR_NO_MATCHING_METHOD {
				ret := NewMessage(TYPE_NONCONFIRMABLE, COAPCODE_405_METHOD_NOT_ALLOWED, msg.MessageId)
				ret.CloneOptions(msg, OPTION_URI_PATH, OPTION_CONTENT_FORMAT)

				SendMessageTo(ret, conn, addr)
				CallEvent(EVT_ERROR, s.events[EVT_ERROR])
				return
			}

			if err == ERR_UNSUPPORTED_CONTENT_FORMAT {
				ret := NewMessage(TYPE_NONCONFIRMABLE, COAPCODE_415_UNSUPPORTED_CONTENT_FORMAT, msg.MessageId)
				ret.CloneOptions(msg, OPTION_URI_PATH, OPTION_CONTENT_FORMAT)

				SendMessageTo(ret, conn, addr)
				CallEvent(EVT_ERROR, s.events[EVT_ERROR])
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


			req := NewClientRequestFromMessage(msg, attrs, conn, addr)
			req.server = s

			if msg.MessageType == TYPE_CONFIRMABLE {
				if msg.GetOption(OPTION_OBSERVE) != nil {
					// TODO: if server doesn't allow observing, return error

					// Observe Request & Fire OnObserve Event
					CallEvent(EVT_OBSERVE, s.events[EVT_OBSERVE])

					// Register observation of client
					s.addObservation(msg.GetUriPath(), string(msg.Token), addr)
					req.Observe()
					req.SetConfirmable(false)

					SendMessageTo(req.GetMessage(), conn, addr)
				}
			} else {
				log.Println("Non Confirmable Message .. check observation")
			}

			resp := route.Handler(req).(*CoapResponse)

			// TODO: Validate Message before sending (.e.g missing messageId)
			SendMessageTo(resp.GetMessage(), conn, addr)
		}
	}
}

func (s *CoapServer) NewRoute(path string, method CoapCode, fn RouteHandler) *Route {
	route := CreateNewRoute(path, MethodString(method), fn)
	s.routes = append(s.routes, route)

	return route
}

func (c *CoapServer) ServerSend(req *CoapRequest) (*CoapResponse, error) {
	return SendMessage(req.GetMessage(), c.localConn)
}

func (c *CoapServer) ServerSendTo(req *CoapRequest, addr *net.UDPAddr) (*CoapResponse, error) {
	return SendMessageTo(req.GetMessage(), c.localConn, addr)
}

func (c *CoapServer) ClientSend(req *CoapRequest) (*CoapResponse, error) {
	return SendMessage(req.GetMessage(), c.remoteConn)
}

func (c *CoapServer) ClientSendTo(req *CoapRequest, addr *net.UDPAddr) (*CoapResponse, error) {
	return SendMessageTo(req.GetMessage(), c.remoteConn, addr)
}

func (c *CoapServer) On(e EventCode, fn FnCanopusEvent) {
	c.events[e] = fn
}

func (c *CoapServer) NotifyChange(resource, value string) {
	t := c.observations[resource]

	if t != nil {
		req := NewRequest(TYPE_CONFIRMABLE, GET, GenerateMessageId())
		req.Observe()

		for _, r := range t {
			req.SetToken(r.Token)
			req.SetStringPayload(value)
			log.Println("NotifyChange", r.Token, req, r.Addr)

			go c.ServerSendTo(req, r.Addr)
		}
	}
}

func (s *CoapServer) addObservation(resource, token string, addr *net.UDPAddr) {
	s.observations[resource] = append(s.observations[resource], NewObservation(addr, token))
}

func (c *CoapServer) Dial(host string) {
	localAddr, err := net.ResolveUDPAddr("udp", ":0")
	log.Println("Local Address", localAddr)
	logging.LogError(err)

	remoteAddr, err := net.ResolveUDPAddr("udp", host)
	logging.LogError(err)

	remoteConn, err := net.DialUDP("udp", localAddr, remoteAddr)
	logging.LogError(err)

	c.remoteAddr = remoteAddr
	c.remoteConn = remoteConn
}

////////////////////////////////////////////////////////////////////////////////
func NewObservation (addr *net.UDPAddr, token string) *Observation {
	return &Observation{
		Addr: addr,
		Token: token,
	}
}

type Observation struct {
	Addr 	*net.UDPAddr
	Token 	string
}
