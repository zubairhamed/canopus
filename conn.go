package goap

import "net"

/*
SendInternalServerError()
SendAcknowledge()
SendConfirmable()
SendNonConfirmable()
SendContent()
*/

func SendError402BadOption(messageId uint16, conn *net.UDPConn, addr *net.UDPAddr) {
	msg := NewMessageOfType(COAPCODE_402_BAD_OPTION, messageId)
	msg.SetStringPayload("Bad Option: An unknown option of type critical was encountered")

	SendMessage(msg, conn, addr)
}

// Sends a CoAP Message to UDP address
func SendMessage(msg *Message, conn *net.UDPConn, addr *net.UDPAddr) error {
	b, _ := MessageToBytes(msg)
	_, err := conn.WriteTo(b, addr)

	return err
}
