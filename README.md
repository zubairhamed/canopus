# Canopus

[![GoDoc](https://godoc.org/github.com/zubairhamed/canopus?status.svg)](https://godoc.org/github.com/zubairhamed/canopus)
[![Build Status](https://drone.io/github.com/zubairhamed/canopus/status.png)](https://drone.io/github.com/zubairhamed/canopus/latest)
[![Coverage Status](https://coveralls.io/repos/zubairhamed/canopus/badge.svg?branch=master)](https://coveralls.io/r/zubairhamed/canopus?branch=master)

#### Canopus is a client/server implementation of the [Constrained Application Protocol (CoAP)][RFC7252]

[RFC7252]: http://tools.ietf.org/html/rfc7252

### Example
```go
package main

import (
	. "github.com/zubairhamed/canopus"
	"log"
)

func main() {
	client := NewCoapServer(":0")

	client.OnStart(func (server *CoapServer){
		client.Dial("localhost:5683")

		req := NewRequest(TYPE_CONFIRMABLE, GET, 50782)
		req.SetStringPayload("Hello, canopus")
		req.SetRequestURI("/hello")

		resp, err := client.Send(req)
		if err != nil {
			log.Println(err)
		} else {
			log.Println("Got Synchronous Response:")
			log.Println(CoapCodeToString(resp.GetMessage().Code))
		}
	})

	client.Start()
}
```

### Server Example
```go
	server := NewLocalServer()

	server.NewRoute("hello", GET, func(r network.Request) network.Response {
		req := r.(*CoapRequest)
		msg := NewMessageOfType(TYPE_ACKNOWLEDGEMENT, req.GetMessage().MessageId)
		msg.SetStringPayload("Acknowledged")
		res := NewResponse(msg, nil)
                                      
		return res
	})

	server.Start()
```

### Observe / Notify

#### Server
```go
package main
import (
	. "github.com/zubairhamed/canopus"
	"github.com/zubairhamed/go-commons/network"
	"time"
	"log"
)

func main() {
	server := NewLocalServer()
	server.NewRoute("watch/this", GET, routeHandler)

	GenerateRandomChangeNotifications(server)

	server.OnMessage(func (msg *Message, inbound bool){
		// PrintMessage(msg)
	})

	server.OnObserve(func(resource string, msg *Message){
		log.Println("Observe Requested for " + resource)
	})

	server.Start()
}

func GenerateRandomChangeNotifications(server *CoapServer) {
	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				log.Println("Notify Change..")
				server.NotifyChange("watch/this", "Some new value", false)
			}
		}
	}()
}

func routeHandler(r network.Request) network.Response {
	req := r.(*CoapRequest)
	msg := NewMessageOfType(TYPE_ACKNOWLEDGEMENT, req.GetMessage().MessageId)
	msg.SetStringPayload("Acknowledged")
	res := NewResponse(msg, nil)

	return res
}
```

#### Client
```go
package main

import (
	. "github.com/zubairhamed/canopus"
	"log"
)

func main() {
	client := NewCoapServer(":0")

	client.OnStart(func (server *CoapServer){
		client.Dial("localhost:5683")
		req := NewRequest(TYPE_CONFIRMABLE, GET, GenerateMessageId())
		req.SetRequestURI("watch/this")
		req.Observe(0)

		_, err := client.Send(req)
		if err != nil {
			log.Println(err)
		}
	})

	client.OnNotify(func (resource string, value interface{}, msg *Message) {
		// PrintMessage(msg)
		log.Println("Got Change Notification for resource and value: ", resource, value)
	})

	client.Start()
}
```