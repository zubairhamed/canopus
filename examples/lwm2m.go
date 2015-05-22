package main
import (
    . "github.com/zubairhamed/goap"
    "log"
    "errors"
)

func main() {
    server := NewLocalServer()

    server.NewRoute("bs", POST, bootstrap)
    server.NewRoute("rd", POST, registration)
    server.NewRoute("go-lwm2m", POST, bootstrap)

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

    ep := req.GetUriQuery("ep")
    lt := req.GetUriQuery("lt")
    b := req.GetUriQuery("b")

    log.Println(ep, lt, b)

    /*
        - Record IP and Port for all future interactions
        - Record registration
        - If already registered, remove and re-register



2015/04/14 09:47:01 Got Synchronous Response:
2015/04/14 09:47:01 201 Created
2015/04/14 09:47:01 = = = = = = = = = = = = = = = =
2015/04/14 09:47:01 Code:  65
2015/04/14 09:47:01 MessageId:  0
2015/04/14 09:47:01 MessageType:  2
2015/04/14 09:47:01 Token:  UQ4kgpiE
2015/04/14 09:47:01 Token Length:  8
2015/04/14 09:47:01 Payload:
2015/04/14 09:47:01  - - - OPTIONS - - -
2015/04/14 09:47:01 Code/Number:  8 , Name:  Location-Path , Value:  rd
2015/04/14 09:47:01 Code/Number:  8 , Name:  Location-Path , Value:  WzP8ECttu3
2015/04/14 09:47:01 = = = = = = = = = = = = = = = =
    */



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

// Errors
var ERR_LWM2M_MANDATORY_PARAM_NOT_SPECIFIED = errors.New("The mandatory parameter is not specified or unknown parameter is specified")
/*
    Unknown Endpoint Client Name
    Endpoint Client Name does not match with CN field of X.509 Certificates

    The Endpoint Client Name results in a duplicate entry on the LWM2M Server.

    The mandatory parameter is not specified or unknown parameter is specified

    URI of “Update” operation is not found

    “De-register” operation is completed successfully

    URI of “De-register” operation is not found

*/