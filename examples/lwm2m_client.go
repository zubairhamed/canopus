package main
import (
    . "github.com/zubairhamed/goap"
    "log"
)

func main() {
    server := NewServer("udp", COAP_DEFAULT_HOST, 5685)

    server.NewRoute("{obj}/",GET ,func (req *CoapRequest) *CoapResponse {
        log.Println(req)

        return nil
    })

    server.NewRoute("{obj}/{obj2}",GET ,func (req *CoapRequest) *CoapResponse {
        log.Println(req)

        return nil
    })

    server.OnStartup(func (evt *Event){
        client := NewClient()
        defer client.Close()

        client.Dial("udp", COAP_DEFAULT_HOST, COAP_DEFAULT_PORT)

        log.Println(" ---- > BOOTSTRAP")
        req := NewRequest(TYPE_CONFIRMABLE, POST, GenerateMessageId())
        req.SetRequestURI("bs")
        req.SetUriQuery("ep", "GOAPLWM2M")
        resp, err := client.Send(req)

        if err != nil {
            log.Println(err)
        } else {
            PrintMessage(resp.GetMessage())
        }

        log.Println(" ---- > REGISTER")
        req = NewRequest(TYPE_CONFIRMABLE, POST, GenerateMessageId())
        req.SetStringPayload("</1>,</2>,</3>,</4>,</5>,</6>,</7>,</8>,</9>,</10>")
        req.SetRequestURI("rd")
        req.SetUriQuery("ep", "GOAPLWM2M")
        resp, err = client.Send(req)

        var path string
        if err != nil {
            log.Println(err)
        } else {
            PrintMessage(resp.GetMessage())

            path = resp.GetMessage().GetLocationPath()
        }

        // Update
        log.Println(" ---- > UPDATE")
        req = NewRequest(TYPE_CONFIRMABLE, PUT, GenerateMessageId())
        req.SetRequestURI(path)
        resp, err = client.Send(req)
        if err != nil {
            log.Println(err)
        } else {
            PrintMessage(resp.GetMessage())
        }

        // Delete
        /*
        log.Println(" ---- > DELETE")
        req = NewRequest(TYPE_CONFIRMABLE, DELETE, GenerateMessageId())
        req.SetRequestURI(path)
        resp, err = client.Send(req)
        if err != nil {
            log.Println(err)
        } else {
            PrintMessage(resp.GetMessage())
        }
        */
    })

    server.Start()
}

/*
    **** Bootstrap ****
    ## Request Bootstrap
    POST
    /bs?ep={Endpoint Client Name}
    > 2.04 Changed
    < 4.00 Bad Request

    ## Write
    PUT
    /{Object ID}/{Object Instance ID}/ {Resource ID}
    > 2.04 Changed
    < 4.00 Bad Request

    ## Delete
    DELETE
    /{Object ID}/{Object Instance ID}
    > 2.02 Deleted
    < 4.05 Method Not Allowed

    **** Registration ****
    /rd
    ## Register
    POST
    rd?ep={Endpoint Client Name}&lt={Lifetime}&sms={MSISDN} &lwm2m={version}&b={binding}
    > 2.01 Created
    < 4.00 Bad Request, 4.09 Conflict

    ## Update
    PUT
    /{location}?lt={Lifetime}&sms={MSISDN} &b={binding}
    > 2.01 Created
    < 4.00 Bad Request, 4.04 Not Found

    ## Delete
    DELETE
    /{location}
    > 2.02 Deleted
    < 4.04 Not Found

    **** Device Management & Service Enablement Interface ****
    ## Read
    GET
    /{Object ID}/{Object Instance ID}/{Resource ID}
    > 2.05 Content
    < 4.01 Unauthorized, 4.04 Not Found, 4.05 Method Not Allowed

    ## Discover
    GET + Accept: application/link- forma
    /{Object ID}/{Object Instance ID}/{Resource ID}
    > 2.05 Content
    < 4.04 Not Found, 4.01 Unauthorized, 4.05 Method Not Allowed

    ## Write
    PUT / POST
    /{Object ID}/{Object Instance ID}/{Resource ID}
    > 2.04 Changed
    < 4.00 Bad Request, 4.04 Not Found, 4.01 Unauthorized, 4.05 Method Not Allowed

    ## Write Attributes
    PUT
    /{Object ID}/{Object Instance ID}/{Resource ID}?pmin={minimum period}&pmax={maximum period}&gt={greater than}&lt={less than}&st={step}&cancel
    > 2.04 Changed
    < 4.00 Bad Request, 4.04 Not Found, 4.01 Unauthorized, 4.05 Method Not Allowed

    ## Execute
    POST
    /{Object ID}/{Object Instance ID}/{Resource ID}
    > 2.04 Changed
    < 4.00 Bad Request, 4.01 Unauthorized, 4.04 Not Found, 4.05 Method Not Allowed

    ## Create
    POST
    /{Object ID}/{Object Instance ID}
    > 2.01 Created
    < 4.00 Bad Request, 4.01 Unauthorized, 4.04 Not Found, 4.05 Method Not Allowed

    ## Delete
    DELETE
    /{Object ID}/{Object Instance ID}
    > 2.02 Deleted
    < 4.01 Unauthorized, 4.04 Not Found, 4.05 Method Not Allowed

    **** Information Reporting Interface ****
    ## Observe
    GET + Observe option
    /{Object ID}/{Object Instance ID}/{Resource ID}
    > 2.05 Content with Observe option
    < 4.04 Not Found, 4.05 Method Not Allowed

    ## Cancel Observe
    Reset message

    ## Notify
    Asynchronous Response
    > 2.04 Changed

*/
