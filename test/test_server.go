package main

import (
    "github.com/zubairhamed/goap"
    "fmt"
)

func main() {
    fmt.Println("GoAP Server Test Started")
    server := goap.NewServer("udp", ":10001")

    server.Handle("example", goap.METHOD_GET , handleGet)
    server.Handle("example", goap.METHOD_DELETE , handleDelete)
    server.Handle("example", goap.METHOD_POST , handlePost)
    server.Handle("example", goap.METHOD_PUT , handlePut)

    server.Start()
}

func handleGet(msg goap.Message) goap.Message {
    fmt.Println("Got GET Message:")
    fmt.Println(goap.PayloadAsString(msg.Payload()))

    return msg
}

func handlePost(msg goap.Message) goap.Message {
    fmt.Println("Got POST Message:")
    fmt.Println(goap.PayloadAsString(msg.Payload()))

    return msg
}

func handlePut(msg goap.Message) goap.Message {
    fmt.Println("Got PUT Message:")
    fmt.Println(goap.PayloadAsString(msg.Payload()))

    return msg

}

func handleDelete(msg goap.Message) goap.Message {
    fmt.Println("Got DELETE Message:")
    fmt.Println(goap.PayloadAsString(msg.Payload()))

    return msg
}