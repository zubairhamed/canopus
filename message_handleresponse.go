package canopus

func handleAcknowledgeObserveRequest(s CoapServer, msg Message) {
	s.GetEvents().Notify(msg.GetURIPath(), msg.GetPayload(), msg)
}
