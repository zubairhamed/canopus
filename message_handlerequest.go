package canopus

import "log"

func handleRequest(s CoapServer, msg *Message, session *Session) {
	if msg.MessageType != MessageReset {
		// Unsupported Method
		if msg.Code != Get && msg.Code != Post && msg.Code != Put && msg.Code != Delete {
			handleReqUnsupportedMethodRequest(msg, session)
			return
		}

		//if err != nil {
		//	s.GetEvents().Error(err)
		//	if err == ErrUnknownCriticalOption {
		//		handleReqUnknownCriticalOption(s, msg, conn, addr)
		//		return
		//	}
		//}

		// Proxy
		if IsProxyRequest(msg) {
			handleReqProxyRequest(s, msg, session)
		} else {
			route, attrs, err := MatchingRoute(msg.GetURIPath(), MethodString(msg.Code), msg.GetOptions(OptionContentFormat), s.GetRoutes())
			if err != nil {
				s.GetEvents().Error(err)
				if err == ErrNoMatchingRoute {
					handleReqNoMatchingRoute(msg, session)
					return
				}

				if err == ErrNoMatchingMethod {
					handleReqNoMatchingMethod(msg, session)
					return
				}

				if err == ErrUnsupportedContentFormat {
					handleReqUnsupportedContentFormat(msg, session)
					return
				}

				log.Println("Error occured parsing inbound message")
				return
			}

			// Duplicate Message ID Check
			if s.IsDuplicateMessage(msg) {
				PrintMessage(msg)
				if msg.MessageType == MessageConfirmable {
					log.Println("Duplicate Message ID ", msg.MessageID)
					handleReqDuplicateMessageID(msg, session)
				}
				return
			}

			s.UpdateMessageTS(msg)

			// Auto acknowledge
			// TODO: Necessary?
			if msg.MessageType == MessageConfirmable && route.AutoAck {
				handleRequestAcknowledge(msg, session)
			}
			req := NewClientRequestFromMessage(msg, attrs, session)
			if msg.MessageType == MessageConfirmable {

				// Observation Request
				obsOpt := msg.GetOption(OptionObserve)
				if obsOpt != nil {
					handleReqObserve(s, req, msg, session)
				}
			}
			opt := req.GetMessage().GetOption(OptionBlock1)
			if opt != nil {
				blockOpt := Block1OptionFromOption(opt)

				// 0000 1 010
				/*
									[NUM][M][SZX]
									2 ^ (2 + 4)
									2 ^ 6 = 32
									Size = 2 ^ (SZX + 4)

									The value 7 for SZX (which would
					      	indicate a block size of 2048) is reserved, i.e. MUST NOT be sent
					      	and MUST lead to a 4.00 Bad Request response code upon reception
					      	in a request.
				*/

				if blockOpt.Value != nil {
					if blockOpt.Code == OptionBlock1 {
						exp := blockOpt.Exponent()

						if exp == 7 {
							handleReqBadRequest(msg, session)
							return
						}

						// szx := blockOpt.Size()
						hasMore := blockOpt.HasMore()
						seqNum := blockOpt.Sequence()
						// fmt.Println("Out Values == ", blockOpt.Value, exp, szx, 2, hasMore, seqNum)

						s.GetEvents().BlockMessage(msg, true)

						s.UpdateBlockMessageFragment(session.GetAddress().String(), msg, seqNum)

						if hasMore {
							handleReqContinue(msg, session)
							// Auto Respond to client

						} else {
							// TODO: Check if message is too large
							msg = NewMessage(msg.MessageType, msg.Code, msg.MessageID)
							msg.Payload = s.FlushBlockMessagePayload(session.GetAddress().String())
							req = NewClientRequestFromMessage(msg, attrs, session)
						}
					} else if blockOpt.Code == OptionBlock2 {

					} else {
						// TOOO: Invalid Block option Code
					}
				}
			}
			resp := route.Handler(req)
			_, nilresponse := resp.(NilResponse)
			if !nilresponse {
				respMsg := resp.GetMessage()
				respMsg.Token = req.GetMessage().Token

				// TODO: Validate Message before sending (e.g missing messageId)
				err := ValidateMessage(respMsg)
				if err == nil {
					s.GetEvents().Message(respMsg, false)
					SendMessage(respMsg, session)
				}
			}
		}
	}
}

func handleReqUnknownCriticalOption(msg *Message, session *Session) {
	if msg.MessageType == MessageConfirmable {
		SendMessage(BadOptionMessage(msg.MessageID, MessageAcknowledgment), session)
	}
	return
}

func handleReqBadRequest(msg *Message, session *Session) {
	if msg.MessageType == MessageConfirmable {
		SendMessage(BadRequestMessage(msg.MessageID, msg.MessageType), session)
	}
	return
}

func handleReqContinue(msg *Message, session *Session) {
	if msg.MessageType == MessageConfirmable {
		SendMessage(ContinueMessage(msg.MessageID, msg.MessageType), session)
	}
	return
}

func handleReqUnsupportedMethodRequest(msg *Message, session *Session) {
	ret := NotImplementedMessage(msg.MessageID, MessageAcknowledgment)
	ret.CloneOptions(msg, OptionURIPath, OptionContentFormat)

	// c.GetEvents().Message(ret, false)
	SendMessage(ret, session)
}

func handleReqProxyRequest(s CoapServer, msg *Message, session *Session) {
	if !s.AllowProxyForwarding(msg, session.GetAddress()) {
		SendMessage(ForbiddenMessage(msg.MessageID, MessageAcknowledgment), session)
	}

	proxyURI := msg.GetOption(OptionProxyURI).StringValue()
	if IsCoapURI(proxyURI) {
		s.ForwardCoap(msg, session)
	} else if IsHTTPURI(proxyURI) {
		s.ForwardHTTP(msg, session)
	} else {
		//
	}
}

func handleReqNoMatchingRoute(msg *Message, session *Session) {
	ret := NotFoundMessage(msg.MessageID, MessageAcknowledgment)
	ret.CloneOptions(msg, OptionURIPath, OptionContentFormat)
	ret.Token = msg.Token

	SendMessage(ret, session)
}

func handleReqNoMatchingMethod(msg *Message, session *Session) {
	ret := MethodNotAllowedMessage(msg.MessageID, MessageAcknowledgment)
	ret.CloneOptions(msg, OptionURIPath, OptionContentFormat)

	SendMessage(ret, session)
}

func handleReqUnsupportedContentFormat(msg *Message, session *Session) {
	ret := UnsupportedContentFormatMessage(msg.MessageID, MessageAcknowledgment)
	ret.CloneOptions(msg, OptionURIPath, OptionContentFormat)

	// s.GetEvents().Message(ret, false)
	SendMessage(ret, session)
}

func handleReqDuplicateMessageID(msg *Message, session *Session) {
	ret := EmptyMessage(msg.MessageID, MessageReset)
	ret.CloneOptions(msg, OptionURIPath, OptionContentFormat)

	SendMessage(ret, session)
}

func handleRequestAcknowledge(msg *Message, session *Session) {
	ack := NewMessageOfType(MessageAcknowledgment, msg.MessageID)

	SendMessage(ack, session)
}

func handleReqObserve(s CoapServer, req CoapRequest, msg *Message, session *Session) {
	// TODO: if server doesn't allow observing, return error
	addr := session.GetAddress()

	// TODO: Check if observation has been registered, if yes, remove it (observation == cancel)
	resource := msg.GetURIPath()
	if s.HasObservation(resource, addr) {
		// Remove observation of client
		s.RemoveObservation(resource, addr)

		// Observe Cancel Request & Fire OnObserveCancel Event
		s.GetEvents().ObserveCancelled(resource, msg)
	} else {
		// Register observation of client
		s.AddObservation(msg.GetURIPath(), string(msg.Token), session)

		// Observe Request & Fire OnObserve Event
		s.GetEvents().Observe(resource, msg)
	}

	req.GetMessage().AddOption(OptionObserve, 1)
}
