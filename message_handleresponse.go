package canopus

import (
	"log"
	"net"
)

func handleResponse(s CoapServer, msg *Message, conn *net.UDPConn, addr *net.UDPAddr) {
	if msg.GetOption(OptionObserve) != nil {
		handleAcknowledgeObserveRequest(s, msg)
		return
	}

	log.Println("Incoming Response", msg.MessageID)
	ch := GetResponseChannel(s, msg.MessageID)
	if ch != nil {
		PrintMessage(msg)
		resp := &CoapResponseChannel{
			Response: NewResponse(msg, nil),
		}
		ch <- resp
		DeleteResponseChannel(s, msg.MessageID)
	} else {
		log.Println("Channel is nil", msg.MessageID)
	}
}

func handleAcknowledgeObserveRequest(s CoapServer, msg *Message) {
	s.GetEvents().Notify(msg.GetURIPath(), msg.Payload, msg)
}
