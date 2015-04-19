package goap

import (
	"bytes"
	"log"
	"net"
	"strconv"
	"time"
)

// Server
func NewServer(host string) *Server {
	addr, _ := net.ResolveUDPAddr("udp", host)
	s := &Server{addr: addr}

	return s
}

func NewLocalServer() *Server {
	return NewServer(":5683")
}

type Server struct {
	addr          	*net.UDPAddr
	messageIds    	map[uint16]time.Time
	routes        	[]*Route
	conn 		  	*net.UDPConn

	evtServerStart	EventHandler
	evtServerError	EventHandler
}

func (s *Server) Send(req *CoapRequest) {

}

func (s *Server) Start() {

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
		ack.Payload = []byte(buf.String())

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
	serveServer(s)
}

func (s *Server) OnStartup(fn EventHandler) {
	s.evtServerStart = fn
}

func (s *Server) OnError(fn EventHandler) {
	s.evtServerError = fn
}


func startServer(s *Server) (*net.UDPConn) {
	s.messageIds = make(map[uint16]time.Time)

	// udpAddr, err := net.ResolveUDPAddr("udp", s.host)
	/*
	if err != nil {
		log.Fatal(err)
	}
	*/

	conn, err := net.ListenUDP("udp", s.addr)
	if err != nil {
		log.Fatal(err)
	}

    log.Println("Started server ", s.conn.LocalAddr())

	CallEvent(s.evtServerStart)

	return conn
}

func handleMessageIdPurge(s *Server) {
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

func serveServer(s *Server) {
	conn := startServer(s)

	handleMessageIdPurge(s)

	s.conn = conn

	readBuf := make([]byte, BUF_SIZE)
	for {
		len, addr, err := conn.ReadFromUDP(readBuf)
		if err == nil {
			log.Println(readBuf)

			msgBuf := make([]byte, len)
			copy(msgBuf, readBuf)

			// Look for route handler matching path and then dispatch
			go s.handleMessage(msgBuf, conn, addr)
		}
	}
}

func (s *Server) handleMessage(msgBuf []byte, conn *net.UDPConn, addr *net.UDPAddr) {
	msg, err := BytesToMessage(msgBuf)

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
			CallEvent(s.evtServerError)
			return
		}

		if err == ERR_NO_MATCHING_METHOD {
			ret := NewMessage(TYPE_NONCONFIRMABLE, COAPCODE_405_METHOD_NOT_ALLOWED, msg.MessageId)
			ret.CloneOptions(msg, OPTION_URI_PATH, OPTION_CONTENT_FORMAT)

            SendMessageTo(ret, conn, addr)
			CallEvent(s.evtServerError)
			return
		}

		if err == ERR_UNSUPPORTED_CONTENT_FORMAT {
			ret := NewMessage(TYPE_NONCONFIRMABLE, COAPCODE_415_UNSUPPORTED_CONTENT_FORMAT, msg.MessageId)
			ret.CloneOptions(msg, OPTION_URI_PATH, OPTION_CONTENT_FORMAT)

            SendMessageTo(ret, conn, addr)
			CallEvent(s.evtServerError)
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

func (s *Server) Close() {
	s.conn.Close()
}
