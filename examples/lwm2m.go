package main
import (
    . "github.com/zubairhamed/goap"
    "log"
)

func main() {
    server := NewLocalServer()

    server.NewRoute("bs", POST, bootstrap)
    server.NewRoute("rd", POST, registration)

    /*
    server.NewRoute("{obj}/{inst}/{rsrc}", PUT, svc)

    server.NewRoute("{obj}/{inst}/{rsrc}", GET, svc).BindMediaTypes(MEDIATYPE_APPLICATION_LINK_FORMAT)
    server.NewRoute("{obj}/{inst}/{rsrc}", GET, svc)
    server.NewRoute("{obj}/{inst}/{rsrc}", GET, svc)    // Get + Observe Options?

    server.NewRoute("{obj}/{inst}/{rsrc}", PUT, svc)
    server.NewRoute("{obj}/{inst}/{rsrc}", POST, svc)

    server.NewRoute("{obj}/{inst}/{rsrc}", GET, svc)    // Reset??
    server.NewRoute("{obj}/{inst}/{rsrc}", GET, svc)    // Notify??

    server.NewRoute("{obj}/{inst}", DELETE, svc)
    server.NewRoute("{obj}/{inst}", POST, svc)

    server.NewRoute("{obj}", PUT, svc)
    server.NewRoute("{obj}", DELETE, svc)
    */

    server.Start()
}

func bootstrap (req *CoapRequest) *CoapResponse {
    return nil
}

func registration (req *CoapRequest) *CoapResponse {
    log.Println("Registration")

    msg := NewMessageOfType(COAPCODE_201_CREATED, req.GetMessage().MessageId)
    resp := NewResponse(msg, nil)

    PrintMessage(req.GetMessage())

    return resp
}

/*
/bs 											POST
/rd												POST

/{obj}/{inst}/{rsrc}			PUT				// Write & Write Attribute

/{obj}/{inst}/{rsrc}			GET
/{obj}/{inst}/{rsrc}			GET 			application/link-format
/{obj}/{inst}/{rsrc}			GET + Observe Options		// Observe

/{obj}/{inst}/{rsrc}			PUT/POST	// Write and Execute


/{obj}/{inst}/{rsrc}			Reset		// Cancel Observation
/{obj}/{inst}/{rsrc}			Async Response	// Notify

/{obj}/{inst}							DELETE		// Delete
/{obj}/{inst}							POST			// Create

/{loc}										PUT
/{loc}										DELETE

*/

// TODO: Service Instantiation for LWM2M and iPSO SmartObjects

//