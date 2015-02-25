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
	ack := goap.DefaultMessage()
	ack.MessageType = goap.TYPE_ACKNOWLEDGEMENT
	ack.MessageId = msg.GetMessageId()
	ack.Payload = []byte("Hello CoAP!")
	ack.Token = msg.GetToken()

    return ack
}

func handlePost(msg goap.Message) goap.Message {
    fmt.Println(goap.PayloadAsString(msg.GetPayload()))

    return msg
}

func handlePut(msg goap.Message) goap.Message {
    fmt.Println("Got PUT Message:")
    fmt.Println(goap.PayloadAsString(msg.GetPayload()))

    return msg

}

func handleDelete(msg goap.Message) goap.Message {
    fmt.Println("Got DELETE Message:")
    fmt.Println(goap.PayloadAsString(msg.GetPayload()))

    return msg
}
