package canopus

func handleResponse(s CoapServer, msg Message, session Session) {
	if msg.GetOption(OptionObserve) != nil {
		handleAcknowledgeObserveRequest(s, msg)
		return
	}

	ch := GetResponseChannel(s, msg.GetMessageId())
	if ch != nil {
		resp := &CoapResponseChannel{
			Response: NewResponse(msg, nil),
		}
		ch <- resp
		DeleteResponseChannel(s, msg.GetMessageId())
	}
}

func handleAcknowledgeObserveRequest(s CoapServer, msg Message) {
	s.GetEvents().Notify(msg.GetURIPath(), msg.GetPayload(), msg)
}
