package canopus

import (
	"fmt"
	"log"
	"net"
	"strconv"
)

func handleRequest(s CoapServer, err error, msg *Message, conn *net.UDPConn, addr *net.UDPAddr) {
	if msg.MessageType != MessageReset {
		// Unsupported Method
		if msg.Code != Get && msg.Code != Post && msg.Code != Put && msg.Code != Delete {
			handleReqUnsupportedMethodRequest(s, msg, conn, addr)
			return
		}

		if err != nil {
			s.GetEvents().Error(err)
			if err == ErrUnknownCriticalOption {
				handleReqUnknownCriticalOption(msg, conn, addr)
				return
			}
		}

		// Proxy
		if IsProxyRequest(msg) {
			handleReqProxyRequest(s, msg, conn, addr)
		} else {
			route, attrs, err := MatchingRoute(msg.GetURIPath(), MethodString(msg.Code), msg.GetOptions(OptionContentFormat), s.GetRoutes())
			if err != nil {
				s.GetEvents().Error(err)
				if err == ErrNoMatchingRoute {
					handleReqNoMatchingRoute(s, msg, conn, addr)
					return
				}

				if err == ErrNoMatchingMethod {
					handleReqNoMatchingMethod(s, msg, conn, addr)
					return
				}

				if err == ErrUnsupportedContentFormat {
					handleReqUnsupportedContentFormat(s, msg, conn, addr)
					return
				}

				log.Println("Error occured parsing inbound message")
				return
			}

			// Duplicate Message ID Check
			if s.IsDuplicateMessage(msg) {
				log.Println("Duplicate Message ID ", msg.MessageID)
				if msg.MessageType == MessageConfirmable {
					handleReqDuplicateMessageID(s, msg, conn, addr)
				}
				return
			}

			s.UpdateMessageTS(msg)

			// Auto acknowledge
			// TODO: Necessary?
			if msg.MessageType == MessageConfirmable && route.AutoAck {
				handleRequestAutoAcknowledge(s, msg, conn, addr)
			}

			req := NewClientRequestFromMessage(msg, attrs, conn, addr)
			if msg.MessageType == MessageConfirmable {

				// Observation Request
				obsOpt := msg.GetOption(OptionObserve)
				if obsOpt != nil {
					handleReqObserve(s, req, msg, conn, addr)
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

				/*
					Server get Block1
						Server saves and append buffer [Ip][BlockXFERs]
						Call OnBlock1 Callback and pass buffer data
						Server Responds OK

						if MORE
							Notify User Callback
							Flush Buffer
						Else
							Do nothing


				*/

				if blockOpt.Value != nil {
					log.Println("HandleRequestBlock1 ## 1")
					if blockOpt.Code == OptionBlock1 {
						log.Println("HandleRequestBlock1 ## 2")
						exp := blockOpt.Exponent()

						if exp == 7 {
							log.Println("HandleRequestBlock1 ## 3")
							handleReqBadRequest(msg, conn, addr)
							return
						}

						szx := blockOpt.Size()
						hasMore := blockOpt.HasMore()
						seqNum := blockOpt.Sequence()
						fmt.Println("Out Values == ", blockOpt.Value, exp, szx, strconv.FormatInt(blockOpt.Value.(int64), 2), hasMore, seqNum)

						s.GetEvents().BlockMessage(msg, true)

						s.UpdateBlockMessageFragment(addr.String(), msg, seqNum)

						if hasMore {
							handleReqContinue(msg, conn, addr)
							// Auto Respond to client
						} else {
							// TODO: Check if message is too large
							msg = NewMessage(msg.MessageType, msg.Code, msg.MessageID)
							msg.Payload = s.FlushBlockMessagePayload(addr.String())
							req = NewClientRequestFromMessage(msg, attrs, conn, addr)
						}
					} else if blockOpt.Code == OptionBlock2 {

					} else {
						// TOOO: Invalid Block option Code
					}

					log.Println("Block Option==", blockOpt, blockOpt.Value)
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

					SendMessageTo(respMsg, NewUDPConnection(conn), addr)
				} else {

				}
			}
		}
	}
}

func handleReqUnknownCriticalOption(msg *Message, conn *net.UDPConn, addr *net.UDPAddr) {
	if msg.MessageType == MessageConfirmable {
		SendMessageTo(BadOptionMessage(msg.MessageID, MessageAcknowledgment), NewUDPConnection(conn), addr)
	}
	return
}

func handleReqBadRequest(msg *Message, conn *net.UDPConn, addr *net.UDPAddr) {
	if msg.MessageType == MessageConfirmable {
		SendMessageTo(BadRequestMessage(msg.MessageID, msg.MessageType), NewUDPConnection(conn), addr)
	}
	return
}

func handleReqContinue(msg *Message, conn *net.UDPConn, addr *net.UDPAddr) {
	if msg.MessageType == MessageConfirmable {
		SendMessageTo(ContinueMessage(msg.MessageID, msg.MessageType), NewUDPConnection(conn), addr)
	}
	return
}

func handleReqUnsupportedMethodRequest(s CoapServer, msg *Message, conn *net.UDPConn, addr *net.UDPAddr) {
	ret := NotImplementedMessage(msg.MessageID, MessageAcknowledgment)
	ret.CloneOptions(msg, OptionURIPath, OptionContentFormat)

	s.GetEvents().Message(ret, false)
	SendMessageTo(ret, NewUDPConnection(conn), addr)
}

func handleReqProxyRequest(s CoapServer, msg *Message, conn *net.UDPConn, addr *net.UDPAddr) {
	if !s.AllowProxyForwarding(msg, addr) {
		SendMessageTo(ForbiddenMessage(msg.MessageID, MessageAcknowledgment), NewUDPConnection(conn), addr)
	}

	proxyURI := msg.GetOption(OptionProxyURI).StringValue()
	if IsCoapURI(proxyURI) {
		s.ForwardCoap(msg, conn, addr)
	} else if IsHTTPURI(proxyURI) {
		s.ForwardHTTP(msg, conn, addr)
	} else {
		//
	}
}

func handleReqNoMatchingRoute(s CoapServer, msg *Message, conn *net.UDPConn, addr *net.UDPAddr) {
	ret := NotFoundMessage(msg.MessageID, MessageAcknowledgment)
	ret.CloneOptions(msg, OptionURIPath, OptionContentFormat)
	ret.Token = msg.Token

	SendMessageTo(ret, NewUDPConnection(conn), addr)
}

func handleReqNoMatchingMethod(s CoapServer, msg *Message, conn *net.UDPConn, addr *net.UDPAddr) {
	ret := MethodNotAllowedMessage(msg.MessageID, MessageAcknowledgment)
	ret.CloneOptions(msg, OptionURIPath, OptionContentFormat)

	s.GetEvents().Message(ret, false)
	SendMessageTo(ret, NewUDPConnection(conn), addr)
}

func handleReqUnsupportedContentFormat(s CoapServer, msg *Message, conn *net.UDPConn, addr *net.UDPAddr) {
	ret := UnsupportedContentFormatMessage(msg.MessageID, MessageAcknowledgment)
	ret.CloneOptions(msg, OptionURIPath, OptionContentFormat)

	s.GetEvents().Message(ret, false)
	SendMessageTo(ret, NewUDPConnection(conn), addr)
}

func handleReqDuplicateMessageID(s CoapServer, msg *Message, conn *net.UDPConn, addr *net.UDPAddr) {
	ret := EmptyMessage(msg.MessageID, MessageReset)
	ret.CloneOptions(msg, OptionURIPath, OptionContentFormat)

	s.GetEvents().Message(ret, false)
	SendMessageTo(ret, NewUDPConnection(conn), addr)
}

func handleRequestAutoAcknowledge(s CoapServer, msg *Message, conn *net.UDPConn, addr *net.UDPAddr) {
	ack := NewMessageOfType(MessageAcknowledgment, msg.MessageID)

	s.GetEvents().Message(ack, false)
	SendMessageTo(ack, NewUDPConnection(conn), addr)
}

func handleReqObserve(s CoapServer, req CoapRequest, msg *Message, conn *net.UDPConn, addr *net.UDPAddr) {
	// TODO: if server doesn't allow observing, return error

	// TODO: Check if observation has been registered, if yes, remove it (observation == cancel)
	resource := msg.GetURIPath()
	if s.HasObservation(resource, addr) {
		// Remove observation of client
		s.RemoveObservation(resource, addr)

		// Observe Cancel Request & Fire OnObserveCancel Event
		s.GetEvents().ObserveCancelled(resource, msg)
	} else {
		// Register observation of client
		s.AddObservation(msg.GetURIPath(), string(msg.Token), addr)

		// Observe Request & Fire OnObserve Event
		s.GetEvents().Observe(resource, msg)
	}

	req.GetMessage().AddOption(OptionObserve, 1)
}
