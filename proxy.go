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
	SendMessageTo(ProxyingNotSupportedMessage(msg.MessageId), NewCanopusUDPConnection(conn), addr)
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

	req, _ := http.NewRequest(MethodString(CoapCode(msg.GetMethod())), proxyUri, nil)

	etag := msg.GetOption(OPTION_ETAG)
	if etag != nil {
		req.Header.Add("ETag", etag.StringValue())
	}

	// TODO: Set timeout handler, and on timeout return 5.04
	resp, err := client.Do(req)

	// TODO: if response not understood or error, return 5.02
	defer resp.Body.Close()

	if err != nil {
		log.Println(err)
		SendMessageTo(BadGatewayMessage(msg.MessageId), NewCanopusUDPConnection(conn), addr)
	}

	contents, _ := ioutil.ReadAll(resp.Body)
	msg.Payload = NewBytesPayload(contents)
	respMsg := NewRequestFromMessage(msg)

	if requestMethod == GET {
		etag := resp.Header.Get("ETag")
		if etag != "" {
			msg.AddOption(OPTION_ETAG, etag)
		}
	}

	SendMessageTo(respMsg.GetMessage(), NewCanopusUDPConnection(conn), addr)
}

// Handles requests for proxying from HTTP to CoAP
func HttpCoapProxyHandler(msg *Message, conn *net.UDPConn, addr *net.UDPAddr) {
	log.Println("HttpCoapProxyHandler Proxy Handler")
}
