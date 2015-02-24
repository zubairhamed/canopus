package goap

import (
    "fmt"
    "net"
)

// Server
func NewServer(net string, host string) Server {
    s := &GoApServer{ net: net, host: host }

    return s
}

type Server interface {
    Handle (path string, method string, fn RouteHandler)
    Start() error
}

type GoApServer struct {
    net     string
    host    string
    routes  map[string] map[Method] RouteHandler
}

func (s *GoApServer) matchingRoute(path string, method Method) (RouteHandler, error) {
    r := s.routes[path]

    if r != nil {
        h := r[method]
        if h != nil {
            return h, nil
        }
    }
    return nil, nil
}

func (s *GoApServer) Handle (path string, method string, fn RouteHandler) {
    fmt.Println("Register Handler")
}

func (s *GoApServer) Start() error {
    fmt.Println("GoAP Server Starting..")

    udpAddr, err := net.ResolveUDPAddr(s.net, s.host);
    if err != nil {
        return err
    }

    conn, err := net.ListenUDP(s.net, udpAddr)
    if err != nil {
        return err
    }

    readBuf := make([]byte, 1500)
    for {
        len, _, err := conn.ReadFromUDP(readBuf)
        if err == nil {
            msgBuf := make([]byte, len)
            copy(msgBuf, readBuf)

            // Look for route handler matching path and then dispatch
            msg := ParseMessage(msgBuf)
            handler, err := s.matchingRoute(msg.Path(), msg.Method())
            if err != nil {
                handler(msg)
            }
        }
    }
}

