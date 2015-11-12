package canopus

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
)

type ProxyHandler func(msg *Message, conn *net.UDPConn, addr *net.UDPAddr)

// The default handler when proxying is disabled
func NullProxyHandler(msg *Message, conn *net.UDPConn, addr *net.UDPAddr) {
	SendMessageTo(ProxyingNotSupportedMessage(msg.MessageId, TYPE_ACKNOWLEDGEMENT), NewCanopusUDPConnection(conn), addr)
}

func CoapCoapProxyHandler(msg *Message, conn *net.UDPConn, addr *net.UDPAddr) {
	/*
		Get value from Proxy-URI
		Resolve Host Address
		Construct CoAP message with Request URI
		Send

		Return response to client

	*/

	log.Println("CoapCoapProxyHandler Proxy Handler")
}

// Handles requests for proxying from CoAP to HTTP
func CoapHttpProxyHandler(msg *Message, conn *net.UDPConn, addr *net.UDPAddr) {
	proxyUri := msg.GetOption(OPTION_PROXY_URI).StringValue()
	requestMethod := msg.Code

	client := &http.Client{}
	req, err := http.NewRequest(MethodString(CoapCode(msg.GetMethod())), proxyUri, nil)
	if err != nil {
		SendMessageTo(BadGatewayMessage(msg.MessageId, TYPE_ACKNOWLEDGEMENT), NewCanopusUDPConnection(conn), addr)
		return
	}

	etag := msg.GetOption(OPTION_ETAG)
	if etag != nil {
		req.Header.Add("ETag", etag.StringValue())
	}

	// TODO: Set timeout handler, and on timeout return 5.04
	resp, err := client.Do(req)
	if err != nil {
		SendMessageTo(BadGatewayMessage(msg.MessageId, TYPE_ACKNOWLEDGEMENT), NewCanopusUDPConnection(conn), addr)
		return
	}

	defer resp.Body.Close()

	contents, _ := ioutil.ReadAll(resp.Body)
	msg.Payload = NewBytesPayload(contents)
	respMsg := NewRequestFromMessage(msg)

	if requestMethod == GET {
		etag := resp.Header.Get("ETag")
		if etag != "" {
			msg.AddOption(OPTION_ETAG, etag)
		}
	}

	// TODO: Check payload length against Size1 options
	if len(respMsg.GetMessage().Payload.String()) > MAX_PACKET_SIZE {
		SendMessageTo(BadGatewayMessage(msg.MessageId, TYPE_ACKNOWLEDGEMENT), NewCanopusUDPConnection(conn), addr)
		return
	}

	_, err = SendMessageTo(respMsg.GetMessage(), NewCanopusUDPConnection(conn), addr)
	if err != nil {
		println(err.Error())
	}
}

// Handles requests for proxying from HTTP to CoAP
func HttpCoapProxyHandler(msg *Message, conn *net.UDPConn, addr *net.UDPAddr) {
	log.Println("HttpCoapProxyHandler Proxy Handler")
}
