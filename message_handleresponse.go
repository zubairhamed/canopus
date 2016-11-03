package canopus

import (
	"log"
	"net"
)

func handleResponse(s CoapServer, msg *Message) {
	if msg.GetOption(OptionObserve) != nil {
		handleAcknowledgeObserveRequest(s, msg)
		return
	}

	ch := GetResponseChannel(s, msg.MessageID)
	if ch != nil {
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
