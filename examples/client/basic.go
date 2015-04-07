package main

import (
	. "github.com/zubairhamed/goap"
	"log"
)

/*	To test this example, also run examples/test_server.go */
func main() {
	// Client
	log.Println("Starting Client..")
	client := NewClient()
	defer client.Close()

	client.Dial("udp", COAP_DEFAULT_HOST, COAP_DEFAULT_PORT)

	msg := NewMessageOfType(TYPE_CONFIRMABLE, 12345)
	msg.Code = GET
	msg.Payload = []byte("Hello, goap")
	msg.AddOptions(NewPathOptions("/example"))
	msg.Token = []byte(GenerateToken(8))

    // Sync Client Test
    log.Println("Sending Synchronous Message")
	resp, err := client.Send(msg)
    if err != nil {
        log.Println(err)
    } else {
        log.Println("Got Synchronous Response:")
        log.Println(resp)
    }

    // Async Client Test
    log.Println("Sending Asynchronous Message")
    client.SendAsync(msg, func(msg *Message, err error){
        if err != nil {
            log.Println(err)
        } else {
            log.Println("Got Asynchronous Response:")
            log.Println(resp)
        }
    })
}
