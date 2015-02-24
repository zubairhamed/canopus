package main

import (
    "github.com/zubairhamed/goap"
    "fmt"
)

func main() {
    fmt.Println("GoAP Server Test Started")
    server := goap.NewServer("udp", ":10001")

    server.Handle("/my/path", "GET", sampleHandler)

    server.Start()
}

func sampleHandler(msg goap.Message) {
    fmt.Println(msg)
}
