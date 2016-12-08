package main

import (
	"fmt"

	"crypto/rand"

	"github.com/zubairhamed/canopus"
)

func main() {
	fmt.Println("Starting up")
	server := canopus.NewServer()

	server.ListenAndServeDTLS(":5682")

	fmt.Println("New Request..")
	req := canopus.NewRequest(canopus.MessageConfirmable, canopus.Post, canopus.GenerateMessageID())
	// req.SetStringPayload(BuildModelResourceStringPayload(c.enabledObjects))
	req.SetRequestURI("/rd")
	req.SetURIQuery("ep", "name")

	// session := canopus.NewUDPServerSession("localhost:5684", server.GetConnection(), server)

	// DTLS Version
	fmt.Println("Initializing SSL")
	secret := make([]byte, 32)
	if n, err := rand.Read(secret); n != 32 || err != nil {
		panic(err)
	}
	server.(*canopus.DefaultCoapServer).SetCookieSecret(secret)
	server.HandlePSK(func(id string) []byte {
		return []byte("secretPSK")
	})

	ctx, err := canopus.NewServerDtlsContext()
	if err != nil {
		panic(err.Error())
	}
	session := canopus.NewDTLSServerSession("localhost:5684", server.GetConnection(), server, ctx)

	fmt.Println("Sending Message")
	resp, err := canopus.SendMessage(req.GetMessage(), session)
	fmt.Println("Response", resp, err)

	canopus.PrintMessage(resp.GetMessage())

	<-make(chan struct{})
}
