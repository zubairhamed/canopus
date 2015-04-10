package main
import (
    . "github.com/zubairhamed/goap"
)

func main() {
    server := NewLocalServer()

    server.NewRoute("{obj}/{inst}/{rsrc}", GET, routeXml)

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

func routeBasic (req *CoapRequest) *CoapResponse {
    msg := NewMessageOfType(TYPE_ACKNOWLEDGEMENT, req.GetMessage().MessageId)
    msg.SetStringPayload("Acknowledged")

    res := NewResponse(msg, nil)

    return res
}

func routeJson (req *CoapRequest) *CoapResponse {
    msg := NewMessageOfType(TYPE_ACKNOWLEDGEMENT, req.GetMessage().MessageId)
    res := NewResponse(msg, nil)

    return res
}

func routeXml (req *CoapRequest) *CoapResponse {
    msg := NewMessageOfType(TYPE_ACKNOWLEDGEMENT, req.GetMessage().MessageId)
    res := NewResponse(msg, nil)

    return res
}

/*
		// goap.PrintMessage(msg)

		fwOpt := msg.GetOption(goap.OPTION_PROXY_URI)
		log.Println(fwOpt)

		ack := goap.NewMessageOfType(goap.TYPE_ACKNOWLEDGEMENT, msg.MessageId)

		return ack

*/
