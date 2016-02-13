package canopus

import "net"

func handleResponse(s CoapServer, msg *Message, conn *net.UDPConn, addr *net.UDPAddr) {
	if msg.GetOption(OptionObserve) != nil {
		handleAcknowledgeObserveRequest(s, msg)
		return
	}
}
func handleAcknowledgeObserveRequest(s CoapServer, msg *Message) {
	s.GetEvents().Notify(msg.GetUriPath(), msg.Payload, msg)
}
