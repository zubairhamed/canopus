package goap

import (
    "fmt"
    "net"
    "errors"
)

// Server
func NewServer(net string, host string) Server {
    s := &GoApServer{ net: net, host: host }
    s.routes = make(map[string] map[uint8] RouteHandler)

    return s
}

type Server interface {
    Handle (path string, method uint8, fn RouteHandler)
    Start() error
}

type GoApServer struct {
    net     string
    host    string
    routes  map[string] map[uint8] RouteHandler
}

func (s *GoApServer) matchingRoute(path string, method uint8) (RouteHandler, error) {
    r := s.routes[path]

    if r != nil {
        h := r[method]
        if h != nil {
            return h, nil
        }
    }
    return nil, errors.New("No matching route found")
}

func (s *GoApServer) Handle (path string, method uint8, fn RouteHandler) {
    if s.routes[path] != nil {
        s.routes[path][method] = fn
    } else {
        m := make(map[uint8] RouteHandler)
        m[method] = fn
        s.routes[path] = m
    }
}

func (s *GoApServer) Start() error {
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
        len, addr, err := conn.ReadFromUDP(readBuf)
        if err == nil {

            msgBuf := make([]byte, len)
            copy(msgBuf, readBuf)

            // Look for route handler matching path and then dispatch
            go s.handleMessage(msgBuf, conn, addr)
        }
    }
}

func (s *GoApServer) handleMessage(msgBuf []byte, conn *net.UDPConn, addr *net.UDPAddr) {
    msg, err := NewMessage(msgBuf)
    if err != nil {
        fmt.Println(err)
        return
    }

    handler, err := s.matchingRoute(msg.Path(), msg.Method())
    if err == nil {
        resp := handler(msg)

        SendPacket (resp, conn, addr)
    }
}

func SendPacket (msg Message, conn *net.UDPConn, addr *net.UDPAddr) error {
    fmt.Println("SendPacket", msg, conn, addr)
    return nil
}
