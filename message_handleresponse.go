package canopus

import "net"

func handleResponse(s *CoapServer, msg *Message, conn *net.UDPConn, addr *net.UDPAddr) {
	if msg.GetOption(OPTION_OBSERVE) != nil {
		handleAcknowledgeObserveRequest(s, msg)
		return
	}
}
func handleAcknowledgeObserveRequest(s *CoapServer, msg *Message) {
	s.events.Notify(msg.GetUriPath(), msg.Payload, msg)
}

