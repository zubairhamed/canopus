package main
import (
    . "github.com/zubairhamed/goap"
    "log"
)

func main() {
    client := NewClient()

    defer client.Close()

    client.Dial("udp", COAP_DEFAULT_HOST, COAP_DEFAULT_PORT)

    req := NewRequest(TYPE_CONFIRMABLE, POST, 12345)
    req.SetStringPayload("<1 />")
    req.SetRequestURI("rd")
    req.SetUriQuery("ep", "DEVKIT")

    resp, err := client.Send(req)
    if err != nil {
        log.Println(err)
    } else {
        log.Println("Got Synchronous Response:")
        log.Println(CoapCodeToString(resp.GetMessage().Code))
        PrintMessage(resp.GetMessage())
    }

    // Update
    req = NewRequest(TYPE_CONFIRMABLE, PUT, 23456)
    req.SetRequestURI("rd")
    resp, err = client.Send(req)
    if err != nil {
        log.Println(err)
    } else {
        log.Println("Got Synchronous Response:")
        log.Println(CoapCodeToString(resp.GetMessage().Code))
        PrintMessage(resp.GetMessage())
    }

}
