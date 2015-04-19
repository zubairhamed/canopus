package main

import (
    . "github.com/zubairhamed/goap"
    "log"
)

/*	To test this example, also run examples/test_server.go */
func main() {
    log.Println("Starting Client..")
    client := NewClient(":64868")
    defer client.Close()

    client.Dial("udp", "127.0.0.1", 58417)

    req := NewRequest(TYPE_CONFIRMABLE, GET, 50782)
    req.SetStringPayload("Hello, GoAP")
    req.SetRequestURI("/0/1/2")

    // Sync Client Test
    log.Println("Sending Synchronous Message")
    resp, err := client.Send(req)
    if err != nil {
        log.Println(err)
    } else {
        log.Println("Got Synchronous Response:")
        log.Println(CoapCodeToString(resp.GetMessage().Code))
    }

    // Async Client Test
    req.IncrementMessageId()
    log.Println("Sending Asynchronous Message")
    client.SendAsync(req, func(resp *CoapResponse, err error){
        if err != nil {
            log.Println(err)
        } else {
            log.Println("Got Asynchronous Response:")
            log.Println(CoapCodeToString(resp.GetMessage().Code))
        }
    })

    // Discovery Test
    /*
    client.Discover(func(resp *CoapCode, err error){

    })
    */
}
