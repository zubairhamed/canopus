package goap

import (
    "fmt"
    "net"
    "errors"
)

const BUF_SIZE = 1500

// Server
func NewServer(net string, host string) *Server {
    s := &Server{ net: net, host: host }
    // s.routes = make(map[string] *Route)

    return s
}

type Server struct {
    net     string
    host    string
	routes	[]*Route
}

func (s *Server) Start() error {
    udpAddr, err := net.ResolveUDPAddr(s.net, s.host);
    if err != nil {
        return err
    }

    conn, err := net.ListenUDP(s.net, udpAddr)
    if err != nil {
        return err
    }

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

    if err == nil {
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
	return &Route{}, errors.New("No matching route found")
}

func (s *Server) NewRoute(path string, fn RouteHandler, method uint8) (*Route) {
	r := &Route{
		AutoAck: true,
		Path: path,
		Method: method,
		Handler: fn,
	}

	s.routes = append(s.routes, r)

	return r
}

func SendPacket (msg *Message, conn *net.UDPConn, addr *net.UDPAddr) error {
	b := MessageToBytes(msg)
	_, err := conn.WriteTo(b, addr)

    return err
}
