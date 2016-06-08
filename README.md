# Canopus

[![GoDoc](https://godoc.org/github.com/zubairhamed/canopus?status.svg)](https://godoc.org/github.com/zubairhamed/canopus)
[![Build Status](https://drone.io/github.com/zubairhamed/canopus/status.png?)](https://drone.io/github.com/zubairhamed/canopus/latest)
[![Coverage Status](https://coveralls.io/repos/zubairhamed/canopus/badge.svg?branch=master)](https://coveralls.io/r/zubairhamed/canopus?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/zubairhamed/canopus)](https://goreportcard.com/report/github.com/zubairhamed/canopus)

#### Canopus is a client/server implementation of the [Constrained Application Protocol (CoAP)][RFC7252]

[RFC7252]: http://tools.ietf.org/html/rfc7252

### Example
```go
package main

import (
	. "github.com/zubairhamed/canopus"
	"log"
)

/*	To test this example, also run examples/test_server.go */
func main() {
	client := NewCoapServer(":0")

	client.OnStart(func (server *CoapServer){
		client.Dial("localhost:5683")

		req := NewRequest(TYPE_CONFIRMABLE, GET, GenerateMessageId())
		req.SetStringPayload("Hello, canopus")
		req.SetRequestURI("/hello")

		resp, err := client.Send(req)
		if err != nil {
			log.Println(err)
		} else {
			log.Println("Got Response:")
			log.Println(resp.GetMessage().Payload.String())
		}
	})

	client.OnMessage(func(msg *Message, inbound bool){
		if inbound {
			log.Println(">>>>> INBOUND <<<<<")
		} else {
			log.Println(">>>>> OUTBOUND <<<<<")
		}

		PrintMessage(msg)
	})

	client.Start()
}
```

### Server Example
```go
	server := NewLocalServer()

	server.Get("/hello", func(req CoapRequest) CoapResponse{
		msg := NewMessageOfType(TYPE_ACKNOWLEDGEMENT, req.GetMessage().MessageId)
		msg.SetStringPayload("Acknowledged: " + req.GetMessage().Payload.String())
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
	"time"
	"log"
	"math/rand"
	"strconv"
)

func main() {
	server := NewLocalServer()
	server.Get("/watch/this", routeHandler)

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
	ticker := time.NewTicker(3 * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				changeVal := strconv.Itoa(rand.Int())
				log.Println("Notify Change..", changeVal)

				server.NotifyChange("/watch/this", changeVal, false)
			}
		}
	}()
}

func routeHandler(req CoapRequest) CoapResponse {
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
		req.SetRequestURI("/watch/this")
		req.GetMessage().AddOption(OPTION_OBSERVE, 0)

		_, err := client.Send(req)
		if err != nil {
			log.Println(err)
		}
	})

	var notifyCount int = 0
	client.OnNotify(func (resource string, value interface{}, msg *Message) {
		if notifyCount < 4 {
			notifyCount++
			log.Println("Got Change Notification for resource and value: ", notifyCount, resource, value)
		} else {
			log.Println("Cancelling Observation after 4 notifications")
			req := NewRequest(TYPE_CONFIRMABLE, GET, GenerateMessageId())
			req.SetRequestURI("watch/this")
			req.GetMessage().AddOption(OPTION_OBSERVE, 0)

			_, err := client.Send(req)
			if err != nil {
				log.Println(err)
			}
		}
	})

	client.Start()
}
```

#### Forward Proxies
```go
package main

import (
	. "github.com/zubairhamed/canopus"
)

func main() {
	server := NewLocalServer()
	server.ProxyCoap(true)  // Forward CoAP Requests
	server.ProxyHttp(true) // Forward HTTP Requests

    // Defaults to NullProxyFilter, if not set.
    // NullProxyFilter allows all requests through (return = true)
	server.SetProxyFilter(func(*Message, *net.UDPAddr) (bool) {
	    // do some checks, e.g. whitelisting etc

	    // allow forwarding, or false to deny
	    return true
	})

	server.Start()
}
```