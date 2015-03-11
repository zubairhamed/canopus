package main

import (
	. "github.com/zubairhamed/goap"
)

func main() {
	server := NewLocalServer()

	server.NewRoute("serviceA", GET, service).BindMediaTypes([]MediaType{MEDIATYPE_APPLICATION_JSON})
	server.NewRoute("serviceB", GET, service).BindMediaTypes([]MediaType{MEDIATYPE_APPLICATION_XML})
	server.NewRoute("serviceC", GET, service).BindMediaTypes([]MediaType{MEDIATYPE_APPLICATION_JSON, MEDIATYPE_TEXT_XML})
	server.NewRoute("serviceD", GET, service)

	server.Start()
}

func service(msg *Message) *Message {
	msg.MessageType = TYPE_ACKNOWLEDGEMENT
	return msg
}
