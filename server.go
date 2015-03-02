package goap

import (
    "fmt"
    "net"
    "errors"
	"log"
    "time"
)

const BUF_SIZE = 1500

const COAP_DEFAULT_HOST = ":5683"

// ERRORS
var ERR_NO_MATCHING_ROUTE = errors.New("No matching route found")


// Server
func NewServer(net string, host string) *Server {
    s := &Server{ net: net, host: host }

    return s
}

type Server struct {
    net         string
    host        string
    messageIds  map[uint16] time.Time
	routes	    []*Route
}

func (s *Server) Start() error {
    s.messageIds = make(map[uint16] time.Time)

    udpAddr, err := net.ResolveUDPAddr(s.net, s.host);
    if err != nil {
        return err
    }

    conn, err := net.ListenUDP(s.net, udpAddr)
    if err != nil {
        return err
    }

    // Routine for clearing up message IDs which has expired
    ticker := time.NewTicker(1 * time.Minute)
    go func() {
        for {
            select {
                case <- ticker.C:
                for k, v := range s.messageIds {
                    elapsed := time.Since(v)
                    if elapsed > 60 {
                        delete(s.messageIds, k)
                    }
                }
            }
        }
    }()

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
    if err != nil {
        fmt.Println(err)
        return
    }

    route, err := s.matchingRoute(msg.GetPath(), msg.Code)
	if err == ERR_NO_MATCHING_ROUTE {
        ret := NewMessageOfType(TYPE_RESET, msg.MessageId)
        ret.Code = COAPCODE_404_NOT_FOUND
        ret.AddOptions(msg.GetOptions(OPTION_URI_PATH))
        ret.AddOptions(msg.GetOptions(OPTION_CONTENT_FORMAT))

		SendPacket(ret, conn, addr)
		return
	}

    // Duplicate Message ID Check
    _, dupe := s.messageIds[msg.MessageId]
    if dupe {
        log.Println("Duplicate Message ID ", msg.MessageId)
        if msg.MessageType == TYPE_CONFIRMABLE {
            ret := NewMessageOfType(TYPE_RESET, msg.MessageId)
            ret.AddOptions(msg.GetOptions(OPTION_URI_PATH))
            ret.AddOptions(msg.GetOptions(OPTION_CONTENT_FORMAT))

            SendPacket(ret, conn, addr)
        }
        return
    }

    if err == nil {
        s.messageIds[msg.MessageId] = time.Now()
		resp := route.Handler(msg)

        SendPacket (resp, conn, addr)
    }
}

func (s *Server) matchingRoute(path string, method uint8) (*Route, error) {
	for _, route := range s.routes {
		if route.Path == path && route.Method == method {
			return route, nil
		}
	}
	return &Route{}, ERR_NO_MATCHING_ROUTE
}

func SendPacket (msg *Message, conn *net.UDPConn, addr *net.UDPAddr) error {
	log.Printf("%#v", msg)

	b := MessageToBytes(msg)
	_, err := conn.WriteTo(b, addr)

    return err
}
