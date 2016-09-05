package canopus

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
)

// Proxy Filter
type ProxyFilter func(*Message, *net.UDPAddr) bool

func NullProxyFilter(*Message, *net.UDPAddr) bool {
	return true
}

type ProxyHandler func(c CoapServer, msg *Message, conn *net.UDPConn, addr *net.UDPAddr)

// The default handler when proxying is disabled
func NullProxyHandler(c CoapServer, msg *Message, conn *net.UDPConn, addr *net.UDPAddr) {
	SendMessageTo(c, ProxyingNotSupportedMessage(msg.MessageID, MessageAcknowledgment), NewUDPConnection(conn), addr)
}

func COAPProxyHandler(c CoapServer, msg *Message, conn *net.UDPConn, addr *net.UDPAddr) {
	proxyURI := msg.GetOption(OptionProxyURI).StringValue()

	parsedURL, err := url.Parse(proxyURI)
	if err != nil {
		log.Println("Error parsing proxy URI")
		SendMessageTo(c, BadGatewayMessage(msg.MessageID, MessageAcknowledgment), NewUDPConnection(conn), addr)
		return
	}

	client := NewCoapClient("Proxy Client")
	client.OnStart(func(server CoapServer) {
		client.Dial(parsedURL.Host)

		msg.RemoveOptions(OptionProxyURI)
		req := NewRequestFromMessage(msg)
		req.SetRequestURI(parsedURL.RequestURI())

		response, err := client.Send(req)
		if err != nil {
			SendMessageTo(c, BadGatewayMessage(msg.MessageID, MessageAcknowledgment), NewUDPConnection(conn), addr)
			client.Stop()
			return
		}

		_, err = SendMessageTo(c, response.GetMessage(), NewUDPConnection(conn), addr)
		if err != nil {
			log.Println("Error occured responding to proxy request")
			client.Stop()
			return
		}
		client.Stop()

	})
	client.Start()
}

// Handles requests for proxying from CoAP to HTTP
func HTTPProxyHandler(c CoapServer, msg *Message, conn *net.UDPConn, addr *net.UDPAddr) {
	proxyURI := msg.GetOption(OptionProxyURI).StringValue()
	requestMethod := msg.Code

	client := &http.Client{}
	req, err := http.NewRequest(MethodString(CoapCode(msg.GetMethod())), proxyURI, nil)
	if err != nil {
		SendMessageTo(c, BadGatewayMessage(msg.MessageID, MessageAcknowledgment), NewUDPConnection(conn), addr)
		return
	}

	etag := msg.GetOption(OptionEtag)
	if etag != nil {
		req.Header.Add("ETag", etag.StringValue())
	}

	// TODO: Set timeout handler, and on timeout return 5.04
	resp, err := client.Do(req)
	if err != nil {
		SendMessageTo(c, BadGatewayMessage(msg.MessageID, MessageAcknowledgment), NewUDPConnection(conn), addr)
		return
	}

	defer resp.Body.Close()

	contents, _ := ioutil.ReadAll(resp.Body)
	msg.Payload = NewBytesPayload(contents)
	respMsg := NewRequestFromMessage(msg)

	if requestMethod == Get {
		etag := resp.Header.Get("ETag")
		if etag != "" {
			msg.AddOption(OptionEtag, etag)
		}
	}

	// TODO: Check payload length against Size1 options
	if len(respMsg.GetMessage().Payload.String()) > MaxPacketSize {
		SendMessageTo(c, BadGatewayMessage(msg.MessageID, MessageAcknowledgment), NewUDPConnection(conn), addr)
		return
	}

	_, err = SendMessageTo(c, respMsg.GetMessage(), NewUDPConnection(conn), addr)
	if err != nil {
		println(err.Error())
	}
}

// Handles requests for proxying from HTTP to CoAP
func HTTPCOAPProxyHandler(msg *Message, conn *net.UDPConn, addr *net.UDPAddr) {
	log.Println("HttpCoapProxyHandler Proxy Handler")
}
