package goap

import (
	"log"
	"net"
)


/*
func GenerateMessageId() uint16 {

}

func GenerateToken() []byte {

}
*/


func DescribeMessage (msg *Message) {
	log.Println("=====================")
	log.Println("Message Type", msg.MessageType)


	log.Println("=================")
	log.Println("==== OPTIONS ====")
	log.Println("=================")

	for _, opt := range msg.Options {
		log.Println(opt.Code, opt.Value)
	}


	log.Println("=====================")
}


func SendPacket (msg *Message, conn *net.UDPConn, addr *net.UDPAddr) error {
	DescribeMessage(msg)

	b := MessageToBytes(msg)
	_, err := conn.WriteTo(b, addr)

	log.Println("Bytes written to conn", b)

	return err
}
