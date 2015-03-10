package main

import (
	"github.com/zubairhamed/goap"
)

func main() {
	server := goap.NewLocalServer()

	server.NewRoute("serviceA", goap.GET, service).BindMediaTypes([]goap.MediaType{goap.MEDIATYPE_APPLICATION_JSON})
	server.NewRoute("serviceB", goap.GET, service).BindMediaTypes([]goap.MediaType{goap.MEDIATYPE_APPLICATION_XML})
	server.NewRoute("serviceC", goap.GET, service).BindMediaTypes([]goap.MediaType{goap.MEDIATYPE_APPLICATION_JSON, goap.MEDIATYPE_TEXT_XML})
	server.NewRoute("serviceD", goap.GET, service)

	server.Start()
}

func service(msg *goap.Message) *goap.Message {
	msg.MessageType = goap.TYPE_ACKNOWLEDGEMENT
	return msg
}
