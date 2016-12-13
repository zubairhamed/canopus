# Canopus

[![GoDoc](https://godoc.org/github.com/zubairhamed/canopus?status.svg)](https://godoc.org/github.com/zubairhamed/canopus)
[![Build Status](https://drone.io/github.com/zubairhamed/canopus/status.png?)](https://drone.io/github.com/zubairhamed/canopus/latest)
[![Coverage Status](https://coveralls.io/repos/zubairhamed/canopus/badge.svg?branch=master)](https://coveralls.io/r/zubairhamed/canopus?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/zubairhamed/canopus)](https://goreportcard.com/report/github.com/zubairhamed/canopus)

#### Canopus is a client/server implementation of the [Constrained Application Protocol (CoAP)][RFC7252]
[RFC7252]: http://tools.ietf.org/html/rfc7252

## Updates
#### 25.11.2016
I've added basic dTLS Support based on [Julien Vermillard's][JVERMILLARD] [implementation][NATIVEDTLS]. Thanks Julien! It should now support PSK-based authentication.
I've also gone ahead and refactored the APIs to make it that bit more Go idiomatic.
[JVERMILLARD]: https://github.com/jvermillard
[NATIVEDTLS]: https://github.com/jvermillard/nativedtls

## Building and running
1. git submodule update --init --recursive
2. cd openssl
3. ./config && make
4. You should then be able to run the examples in the /examples folder

#### Simple Example
```go
	// Server
	// See /examples/simple/server/main.go
	server := canopus.NewServer()

	server.Get("/hello", func(req canopus.Request) canopus.Response {
		msg := canopus.ContentMessage(req.GetMessage().GetMessageId(), canopus.MessageAcknowledgment)
		msg.SetStringPayload("Acknowledged: " + req.GetMessage().GetPayload().String())

		res := canopus.NewResponse(msg, nil)
		return res
	})

	server.ListenAndServe(":5683")

	// Client
	// See /examples/simple/client/main.go
	conn, err := canopus.Dial("localhost:5683")
	if err != nil {
		panic(err.Error())
	}

	req := canopus.NewRequest(canopus.MessageConfirmable, canopus.Get, canopus.GenerateMessageID()).(*canopus.CoapRequest)
	req.SetStringPayload("Hello, canopus")
	req.SetRequestURI("/hello")

	resp, err := conn.Send(req)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Got Response:" + resp.GetMessage().GetPayload().String())
```

#### Observe / Notify
```go
	// Server
	// See /examples/observe/server/main.go
	server := canopus.NewServer()
	server.Get("/watch/this", func(req canopus.Request) canopus.Response {
		msg := canopus.NewMessageOfType(canopus.MessageAcknowledgment, req.GetMessage().GetMessageId(), canopus.NewPlainTextPayload("Acknowledged"))
		res := canopus.NewResponse(msg, nil)

		return res
	})

	ticker := time.NewTicker(3 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				changeVal := strconv.Itoa(rand.Int())
				fmt.Println("[SERVER << ] Change of value -->", changeVal)

				server.NotifyChange("/watch/this", changeVal, false)
			}
		}
	}()

	server.OnObserve(func(resource string, msg canopus.Message) {
		fmt.Println("[SERVER << ] Observe Requested for " + resource)
	})

	server.ListenAndServe(":5683")

	// Client
	// See /examples/observe/client/main.go
	conn, err := canopus.Dial("localhost:5683")

	tok, err := conn.ObserveResource("/watch/this")
	if err != nil {
		panic(err.Error())
	}

	obsChannel := make(chan canopus.ObserveMessage)
	done := make(chan bool)
	go conn.Observe(obsChannel)

	notifyCount := 0
	for {
		select {
		case obsMsg, _ := <-obsChannel:
			if notifyCount == 5 {
				fmt.Println("[CLIENT >> ] Canceling observe after 5 notifications..")
				go conn.CancelObserveResource("watch/this", tok)
				go conn.StopObserve(obsChannel)
				return
			} else {
				notifyCount++
				// msg := obsMsg.Msg\
				resource := obsMsg.GetResource()
				val := obsMsg.GetValue()

				fmt.Println("[CLIENT >> ] Got Change Notification for resource and value: ", notifyCount, resource, val)
			}
		}
	}
```

### dTLS with PSK
```go
	// Server
	// See /examples/dtls/simple-psk/server/main.go
	server := canopus.NewServer()

	server.Get("/hello", func(req canopus.Request) canopus.Response {
		msg := canopus.ContentMessage(req.GetMessage().GetMessageId(), canopus.MessageAcknowledgment)
		msg.SetStringPayload("Acknowledged: " + req.GetMessage().GetPayload().String())
		res := canopus.NewResponse(msg, nil)

		return res
	})

	server.HandlePSK(func(id string) []byte {
		return []byte("secretPSK")
	})

	server.ListenAndServeDTLS(":5684")

	// Client
	// See /examples/dtls/simple-psk/client/main.go
	conn, err := canopus.DialDTLS("localhost:5684", "canopus", "secretPSK")
	if err != nil {
		panic(err.Error())
	}

	req := canopus.NewRequest(canopus.MessageConfirmable, canopus.Get, canopus.GenerateMessageID())
	req.SetStringPayload("Hello, canopus")
	req.SetRequestURI("/hello")

	resp, err := conn.Send(req)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Got Response:" + resp.GetMessage().GetPayload().String())
```

#### CoAP-CoAP Proxy
```go
	// Server
	// See /examples/proxy/coap/server/main.go
	server := canopus.NewServer()

	server.Get("/proxycall", func(req canopus.Request) canopus.Response {
		msg := canopus.ContentMessage(req.GetMessage().GetMessageId(), canopus.MessageAcknowledgment)
		msg.SetStringPayload("Data from :5685 -- " + req.GetMessage().GetPayload().String())
		res := canopus.NewResponse(msg, nil)

		return res
	})
	server.ListenAndServe(":5685")

	// Proxy Server
	// See /examples/proxy/coap/proxy/main.go
	server := canopus.NewServer()
	server.ProxyOverCoap(true)

	server.Get("/proxycall", func(req canopus.Request) canopus.Response {
		canopus.PrintMessage(req.GetMessage())
		msg := canopus.ContentMessage(req.GetMessage().GetMessageId(), canopus.MessageAcknowledgment)
		msg.SetStringPayload("Acknowledged: " + req.GetMessage().GetPayload().String())
		res := canopus.NewResponse(msg, nil)

		return res
	})
	server.ListenAndServe(":5683")

	// Client
	// See /examples/proxy/coap/client/main.go
	conn, err := canopus.Dial("localhost:5683")
	if err != nil {
		panic(err.Error())
	}

	req := canopus.NewRequest(canopus.MessageConfirmable, canopus.Get, canopus.GenerateMessageID())
	req.SetProxyURI("coap://localhost:5685/proxycall")

	resp, err := conn.Send(req)
	if err != nil {
		println("err", err)
	}
	canopus.PrintMessage(resp.GetMessage())
```

#### CoAP-HTTP Proxy
```go
	// Server
	// See /examples/proxy/http/server/main.go
	server := canopus.NewServer()
	server.ProxyOverHttp(true)

	server.ListenAndServe(":5683")

	// Client
	// See /examples/proxy/http/client/main.go
	conn, err := canopus.Dial("localhost:5683")
	if err != nil {
		panic(err.Error())
	}

	req := canopus.NewRequest(canopus.MessageConfirmable, canopus.Get, canopus.GenerateMessageID())
	req.SetProxyURI("https://httpbin.org/get")

	resp, err := conn.Send(req)
	if err != nil {
		println("err", err)
	}
	canopus.PrintMessage(resp.GetMessage())
```
