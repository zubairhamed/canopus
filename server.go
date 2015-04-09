package goap

import (
	"bytes"
	"log"
	"net"
	"strconv"
	"time"
)

// Server
func NewServer(net string, host string, port int) *Server {
	s := &Server{net: net, host: host, port: port, discoveryPort: port}

	return s
}

func NewLocalServer() *Server {
	return NewServer("udp", COAP_DEFAULT_HOST, COAP_DEFAULT_PORT)
}

type Server struct {
	net           string
	host          string
	port          int
	discoveryPort int
	messageIds    map[uint16]time.Time
	routes        []*Route
	conn 		  *net.UDPConn
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

		resp := NewResponseFromMessage(ack)

		return resp
	}

	if s.port == s.discoveryPort {
		s.NewRoute(".well-known/core", GET, discoveryRoute)

		serveServer(s)
	} else {
		discoveryServer := &Server{net: s.net, host: s.host, port: COAP_DEFAULT_PORT}
		discoveryServer.NewRoute(".well-known/core", GET, discoveryRoute)

		serveServer(discoveryServer)
	}
}

func startServer(s *Server) (*net.UDPConn) {
	hostString := s.host + ":" + strconv.Itoa(s.port)
	s.messageIds = make(map[uint16]time.Time)

	udpAddr, err := net.ResolveUDPAddr(s.net, hostString)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.ListenUDP(s.net, udpAddr)
	if err != nil {
		log.Fatal(err)
	}

    log.Println("Started server on port ", s.port)

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

	route, err := MatchingRoute(msg, s.routes)
	if err != nil {
		if err == ERR_NO_MATCHING_ROUTE {
			ret := NewMessage(TYPE_NONCONFIRMABLE, COAPCODE_404_NOT_FOUND, msg.MessageId)
			ret.CloneOptions(msg, OPTION_URI_PATH, OPTION_CONTENT_FORMAT)
			ret.Token = msg.Token

            SendMessageTo(ret, conn, addr)
			return
		}

		if err == ERR_NO_MATCHING_METHOD {
			ret := NewMessage(TYPE_NONCONFIRMABLE, COAPCODE_405_METHOD_NOT_ALLOWED, msg.MessageId)
			ret.CloneOptions(msg, OPTION_URI_PATH, OPTION_CONTENT_FORMAT)

            SendMessageTo(ret, conn, addr)
			return
		}

		if err == ERR_UNSUPPORTED_CONTENT_FORMAT {
			ret := NewMessage(TYPE_NONCONFIRMABLE, COAPCODE_415_UNSUPPORTED_CONTENT_FORMAT, msg.MessageId)
			ret.CloneOptions(msg, OPTION_URI_PATH, OPTION_CONTENT_FORMAT)

            SendMessageTo(ret, conn, addr)
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

		req := NewRequestFromMessage(msg)

		resp := route.Handler(req)

		// TODO: Validate Message before sending (.e.g missing messageId)
        SendMessageTo(resp.GetMessage(), conn, addr)
	}
}

func (s *Server) Close() {
	s.conn.Close()
}
