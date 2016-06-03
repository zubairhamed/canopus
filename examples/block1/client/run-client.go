package main

import (
	"github.com/zubairhamed/canopus"
	"log"
	"io/ioutil"
)

func main() {
	client := canopus.NewCoapServer("0")

	client.OnStart(func(server canopus.CoapServer) {
		client.Dial("localhost:5683")

		file, err := ioutil.ReadFile("./ietf-block.htm")
		if err != nil {
			log.Fatal(err)
		}

		req := canopus.NewRequest(canopus.MessageConfirmable, canopus.Get, canopus.GenerateMessageID())

		blockOpt := canopus.NewBlock1Option(canopus.BlockSize16, true, 0)
		req.GetMessage().SetBlock1Option(blockOpt)
		req.SetPayload(file)
		req.SetRequestURI("/blockinfo")

		// resp, err := client.Send(req)
		_, err = client.Send(req)
		if err != nil {
			log.Println(err)
		} else {
			log.Println("Got Response:")
			// log.Println(resp.GetMessage().Payload.String())
		}
	})

	client.OnError(func(err error) {
		log.Println("An error occured: ", err)
	})

	client.Start()
}
