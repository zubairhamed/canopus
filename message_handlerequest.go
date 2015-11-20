package canopus

import (
	"net"
	"log"
	"time"
)

func handleRequest(s *CoapServer, err error, msg *Message, conn *net.UDPConn, addr *net.UDPAddr) {
	if msg.MessageType != TYPE_RESET {
		// Unsupported Method
		if msg.Code != GET && msg.Code != POST && msg.Code != PUT && msg.Code != DELETE {
			handleReqUnsupportedMethodRequest(s, msg, conn, addr)
			return
		}

		if err != nil {
			s.events.Error(err)
			if err == ERR_UNKNOWN_CRITICAL_OPTION {
				handleReqUnknownCriticalOption(msg, conn, addr)
				return
			}
		}

		// Proxy
		if IsProxyRequest(msg) {
			handleReqProxyRequest(s, msg, conn, addr)
		} else {
			route, attrs, err := MatchingRoute(msg.GetUriPath(), MethodString(msg.Code), msg.GetOptions(OPTION_CONTENT_FORMAT), s.routes)
			if err != nil {
				s.events.Error(err)
				if err == ERR_NO_MATCHING_ROUTE {
					handleReqNoMatchingRoute(s, msg, conn, addr)
					return
				}

				if err == ERR_NO_MATCHING_METHOD {
					handleReqNoMatchingMethod(s, msg, conn, addr)
					return
				}

				if err == ERR_UNSUPPORTED_CONTENT_FORMAT {
					handleReqUnsupportedContentFormat(s, msg, conn, addr)
					return
				}

				log.Println("Error occured parsing inbound message")
				return
			}

			// Duplicate Message ID Check
			_, dupe := s.messageIds[msg.MessageId]
			if dupe {
				log.Println("Duplicate Message ID ", msg.MessageId)
				if msg.MessageType == TYPE_CONFIRMABLE {
					handleReqDuplicateMessageId(s, msg, conn, addr)
				}
				return
			}

			s.messageIds[msg.MessageId] = time.Now()

			// Auto acknowledge
			if msg.MessageType == TYPE_CONFIRMABLE && route.AutoAck {
				handleRequestAutoAcknowledge(s, msg, conn, addr)
			}

			req := NewClientRequestFromMessage(msg, attrs, conn, addr)
			if msg.MessageType == TYPE_CONFIRMABLE {

				// Observation Request
				obsOpt := msg.GetOption(OPTION_OBSERVE)
				if obsOpt != nil {
					handleReqObserve(s, req, msg, conn, addr)
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
					s.events.Message(respMsg, false)

					SendMessageTo(respMsg, NewCanopusUDPConnection(conn), addr)
				} else {

				}
			}
		}
	}
}


func handleReqUnknownCriticalOption(msg *Message, conn *net.UDPConn, addr *net.UDPAddr) {
	if msg.MessageType == TYPE_CONFIRMABLE {
		SendMessageTo(BadOptionMessage(msg.MessageId, TYPE_ACKNOWLEDGEMENT), NewCanopusUDPConnection(conn), addr)
		return
	} else {
		// Ignore silently
		return
	}
}

func handleReqUnsupportedMethodRequest(s *CoapServer, msg *Message, conn *net.UDPConn, addr *net.UDPAddr) {
	ret := NotImplementedMessage(msg.MessageId, TYPE_ACKNOWLEDGEMENT)
	ret.CloneOptions(msg, OPTION_URI_PATH, OPTION_CONTENT_FORMAT)

	s.events.Message(ret, false)
	SendMessageTo(ret, NewCanopusUDPConnection(conn), addr)
}

func handleReqProxyRequest(s *CoapServer, msg *Message, conn *net.UDPConn, addr *net.UDPAddr) {
	if !s.fnProxyFilter(msg, addr) {
		SendMessageTo(ForbiddenMessage(msg.MessageId, TYPE_ACKNOWLEDGEMENT), NewCanopusUDPConnection(conn), addr)
	}

	proxyUri := msg.GetOption(OPTION_PROXY_URI).StringValue()
	if IsCoapUri(proxyUri) {
		s.fnHandleCoapProxy(msg, conn, addr)
	} else if IsHttpUri(proxyUri) {
		s.fnHandleHttpProxy(msg, conn, addr)
	} else {
		//
	}
}

func handleReqNoMatchingRoute(s *CoapServer, msg *Message, conn *net.UDPConn, addr *net.UDPAddr) {
	ret := NotFoundMessage(msg.MessageId, TYPE_ACKNOWLEDGEMENT)
	ret.CloneOptions(msg, OPTION_URI_PATH, OPTION_CONTENT_FORMAT)
	ret.Token = msg.Token

	SendMessageTo(ret, NewCanopusUDPConnection(conn), addr)
}

func handleReqNoMatchingMethod(s *CoapServer, msg *Message, conn *net.UDPConn, addr *net.UDPAddr) {
	ret := MethodNotAllowedMessage(msg.MessageId, TYPE_ACKNOWLEDGEMENT)
	ret.CloneOptions(msg, OPTION_URI_PATH, OPTION_CONTENT_FORMAT)

	s.events.Message(ret, false)
	SendMessageTo(ret, NewCanopusUDPConnection(conn), addr)
}

func handleReqUnsupportedContentFormat(s *CoapServer, msg *Message, conn *net.UDPConn, addr *net.UDPAddr) {
	ret := UnsupportedContentFormatMessage(msg.MessageId, TYPE_ACKNOWLEDGEMENT)
	ret.CloneOptions(msg, OPTION_URI_PATH, OPTION_CONTENT_FORMAT)

	s.events.Message(ret, false)
	SendMessageTo(ret, NewCanopusUDPConnection(conn), addr)
}

func handleReqDuplicateMessageId(s *CoapServer, msg *Message, conn *net.UDPConn, addr *net.UDPAddr) {
	ret := EmptyMessage(msg.MessageId, TYPE_RESET)
	ret.CloneOptions(msg, OPTION_URI_PATH, OPTION_CONTENT_FORMAT)

	s.events.Message(ret, false)
	SendMessageTo(ret, NewCanopusUDPConnection(conn), addr)
}

func handleRequestAutoAcknowledge(s *CoapServer, msg *Message, conn *net.UDPConn, addr *net.UDPAddr) {
	ack := NewMessageOfType(TYPE_ACKNOWLEDGEMENT, msg.MessageId)

	s.events.Message(ack, false)
	SendMessageTo(ack, NewCanopusUDPConnection(conn), addr)
}

func handleReqObserve(s *CoapServer, req CoapRequest, msg *Message, conn *net.UDPConn, addr *net.UDPAddr) {
	// TODO: if server doesn't allow observing, return error

	// TODO: Check if observation has been registered, if yes, remove it (observation == cancel)
	resource := msg.GetUriPath()
	if s.hasObservation(resource, addr) {
		// Remove observation of client
		s.removeObservation(resource, addr)

		// Observe Cancel Request & Fire OnObserveCancel Event
		s.events.ObserveCancelled(resource, msg)
	} else {
		// Register observation of client
		s.addObservation(msg.GetUriPath(), string(msg.Token), addr)

		// Observe Request & Fire OnObserve Event
		s.events.Observe(resource, msg)
	}

	req.GetMessage().AddOption(OPTION_OBSERVE, 1)
}