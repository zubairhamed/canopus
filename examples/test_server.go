package main

import (
    "github.com/zubairhamed/goap"
    "fmt"
	"log"
)

func main() {
    fmt.Println("GoAP Server Test Started")
    server := goap.NewServer("udp", ":10001")

    server.Handle("example", goap.METHOD_GET , handleGet)
    server.Handle("example", goap.METHOD_DELETE , handleDelete)
    server.Handle("example", goap.METHOD_POST , handlePost)
    server.Handle("example", goap.METHOD_PUT , handlePut)

	/*
	server.Handle("example", handleGet).Methods("GET").Schemes("coaps")
	server.Handle("example", handleGet).Methods("GET").Schemes("coaps").AutoAcknowledge().
	*/

    server.Start()
}

func handleGet(msg *goap.Message) *goap.Message {
	ack := goap.NewMessage()
	ack.MessageType = goap.TYPE_ACKNOWLEDGEMENT
	ack.MessageId = msg.MessageId
	ack.Payload = []byte("hello to you!!!!!")
	ack.Token = msg.Token
	ack.Code = goap.COAPCODE_205_CONTENT

	ack.Options = append(ack.Options, goap.NewOption(goap.OPTION_CONTENT_FORMAT, goap.MEDIATYPE_TEXT_PLAIN))

	log.Printf("Transmitting %#v", ack)

    return ack
}

func handlePost(msg *goap.Message) *goap.Message {
    fmt.Println(goap.PayloadAsString(msg.Payload))

    return msg
}

func handlePut(msg *goap.Message) *goap.Message {
    fmt.Println("Got PUT Message:")
    fmt.Println(goap.PayloadAsString(msg.Payload))

    return msg

}

func handleDelete(msg *goap.Message) *goap.Message {
    fmt.Println("Got DELETE Message:")
    fmt.Println(goap.PayloadAsString(msg.Payload))

    return msg
}
