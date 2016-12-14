package main

import (
	"io/ioutil"
	"log"

	"github.com/zubairhamed/canopus"
)

func main() {
	conn, err := canopus.Dial("localhost:5683")

	file, err := ioutil.ReadFile("./ietf-block.htm")
	if err != nil {
		log.Fatal(err)
	}

	req := canopus.NewRequest(canopus.MessageConfirmable, canopus.Post)
	blockOpt := canopus.NewBlock1Option(canopus.BlockSize16, true, 0)

	req.GetMessage().SetBlock1Option(blockOpt)
	req.SetRequestURI("/blockupload")
	req.SetPayload(file)

	resp, err := conn.Send(req)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Got Response:")
		log.Println(resp.GetMessage().GetPayload().String())
	}
}
