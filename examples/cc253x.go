package main

import "github.com/zubairhamed/goap"

func main() {
	server := goap.NewServer("udp", ":10002")
	server.Handle("state", goap.METHOD_POST , addState)
}

func addState(msg *goap.Message) *goap.Message {
	return nil
}
