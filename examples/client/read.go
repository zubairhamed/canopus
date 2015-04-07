package main

import (
    "github.com/zubairhamed/golwm2m"
    "github.com/zubairhamed/goap"
    "log"
)

func main() {
    client := golwm2m.NewClient()

    client.Dial("udp", "localhost", 5683)

    client.Read("1/0", func (msg *goap.Message) {
        log.Println(goap.CoapCodeToString(msg.Code))
    })
}
