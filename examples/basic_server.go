package main

import (
	. "github.com/zubairhamed/canopus"
	"log"
)

func main() {
	server := NewLocalServer()

	server.NewRoute("{obj}/{inst}/{rsrc}", GET, routeParams)

	server.NewRoute("basic", GET, routeBasic)
	server.NewRoute("basic/json", GET, routeJson)
	server.NewRoute("basic/xml", GET, routeXml)

	/*
	   server.OnDiscover(request, response) {

	   }

	   server.OnError(request, error, errorCode) {

	   }
	*/
	server.Start()
}

func routeParams(req *CoapRequest) *CoapResponse {
	msg := NewMessageOfType(TYPE_ACKNOWLEDGEMENT, req.GetMessage().MessageId)
	msg.SetStringPayload("Acknowledged")
	res := NewResponse(msg, nil)

	log.Println(req.GetAttributes())
	log.Println("obj", req.GetAttribute("obj"))
	log.Println("inst", req.GetAttribute("inst"))
	log.Println("rsrc", req.GetAttribute("rsrc"))

	return res
}

func routeBasic(req *CoapRequest) *CoapResponse {
	msg := NewMessageOfType(TYPE_ACKNOWLEDGEMENT, req.GetMessage().MessageId)
	msg.SetStringPayload("Acknowledged")

	res := NewResponse(msg, nil)

	return res
}

func routeJson(req *CoapRequest) *CoapResponse {
	msg := NewMessageOfType(TYPE_ACKNOWLEDGEMENT, req.GetMessage().MessageId)
	res := NewResponse(msg, nil)

	return res
}

func routeXml(req *CoapRequest) *CoapResponse {
	msg := NewMessageOfType(TYPE_ACKNOWLEDGEMENT, req.GetMessage().MessageId)
	res := NewResponse(msg, nil)

	return res
}

/*
	// canopus.PrintMessage(msg)

	fwOpt := msg.GetOption(canopus.OPTION_PROXY_URI)
	log.Println(fwOpt)

	ack := canopus.NewMessageOfType(canopus.TYPE_ACKNOWLEDGEMENT, msg.MessageId)

	return ack

*/
