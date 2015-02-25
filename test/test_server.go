package main

import (
    "github.com/zubairhamed/goap"
    "fmt"
)

func main() {
    fmt.Println("GoAP Server Test Started")
    server := goap.NewServer("udp", ":10001")

    server.Handle("my/path", goap.METHOD_GET , handleGet)
    server.Handle("my/path", goap.METHOD_DELETE , handleDelete)
    server.Handle("my/path", goap.METHOD_POST , handlePost)
    server.Handle("my/path", goap.METHOD_PUT , handlePut)

    server.Start()
}

func handleGet(msg goap.Message) goap.Message {
    fmt.Println("Got GET Message:")
    fmt.Println(msg.Payload())

    return nil
}

func handlePost(msg goap.Message) goap.Message {
    fmt.Println("Got POST Message:")
    fmt.Println(msg.Payload())

    return nil
}

func handlePut(msg goap.Message) goap.Message {
    fmt.Println("Got PUT Message:")
    fmt.Println(msg.Payload())

    return nil

}

func handleDelete(msg goap.Message) goap.Message {
    fmt.Println("Got DELETE Message:")
    fmt.Println(msg.Payload())

    return nil
}