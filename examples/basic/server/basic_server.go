package main

import (
	. "github.com/zubairhamed/canopus"
)

func main() {
	server := NewLocalServer()

//	server.Get("/hello", func (*Request) *CoapResponse {
//
//	})

	server.NewRoute("hello", GET, routeParams)

	server.NewRoute("basic", GET, routeBasic)
	server.NewRoute("basic/json", GET, routeJson)
	server.NewRoute("basic/xml", GET, routeXml)

	server.Start()
}

func routeParams(req *Request) *Response {
	msg := NewMessageOfType(TYPE_ACKNOWLEDGEMENT, req.GetMessage().MessageId)
	msg.SetStringPayload("Acknowledged")
	res := NewResponse(msg, nil)

	return res
}

func routeBasic(req *Request) *Response {
	msg := NewMessageOfType(TYPE_ACKNOWLEDGEMENT, req.GetMessage().MessageId)
	msg.SetStringPayload("Acknowledged")

	res := NewResponse(msg, nil)

	return res
}

func routeJson(req *Request) *Response {
	msg := NewMessageOfType(TYPE_ACKNOWLEDGEMENT, req.GetMessage().MessageId)
	res := NewResponse(msg, nil)

	return res
}

func routeXml(req *Request) *Response {
	msg := NewMessageOfType(TYPE_ACKNOWLEDGEMENT, req.GetMessage().MessageId)
	res := NewResponse(msg, nil)

	return res
}
