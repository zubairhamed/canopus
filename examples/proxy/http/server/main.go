package main

import "github.com/zubairhamed/canopus"

func main() {
	server := canopus.NewServer()
	server.ProxyOverHttp(true)

	server.ListenAndServe(":5683")
	<-make(chan struct{})

}
