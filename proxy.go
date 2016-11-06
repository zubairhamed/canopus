package canopus

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
)

// Proxy Filter
type ProxyFilter func(*Message, net.Addr) bool

func NullProxyFilter(*Message, net.Addr) bool {
	return true
}

type ProxyHandler func(c CoapServer, msg *Message, session Session)

// The default handler when proxying is disabled
func NullProxyHandler(c CoapServer, msg *Message, session Session) {
	SendMessage(ProxyingNotSupportedMessage(msg.MessageID, MessageAcknowledgment), session)
}

func COAPProxyHandler(c CoapServer, msg *Message, session Session) {
	proxyURI := msg.GetOption(OptionProxyURI).StringValue()

	parsedURL, err := url.Parse(proxyURI)
	if err != nil {
		log.Println("Error parsing proxy URI")
		SendMessage(BadGatewayMessage(msg.MessageID, MessageAcknowledgment), session)
		return
	}

	client := NewClient()
	clientConn, err := client.Dial(parsedURL.Host)

	msg.RemoveOptions(OptionProxyURI)
	req := NewRequestFromMessage(msg)
	req.SetRequestURI(parsedURL.RequestURI())

	response, err := clientConn.Send(req)
	if err != nil {
		SendMessage(BadGatewayMessage(msg.MessageID, MessageAcknowledgment), session)
		clientConn.Close()
		return
	}

	_, err = SendMessage(response.GetMessage(), session)
	if err != nil {
		log.Println("Error occured responding to proxy request")
		clientConn.Close()
		return
	}
	clientConn.Close()
}

// Handles requests for proxying from CoAP to HTTP
func HTTPProxyHandler(c CoapServer, msg *Message, session Session) {
	proxyURI := msg.GetOption(OptionProxyURI).StringValue()
	requestMethod := msg.Code

	client := &http.Client{}
	req, err := http.NewRequest(MethodString(CoapCode(msg.GetMethod())), proxyURI, nil)
	if err != nil {
		SendMessage(BadGatewayMessage(msg.MessageID, MessageAcknowledgment), session)
		return
	}

	etag := msg.GetOption(OptionEtag)
	if etag != nil {
		req.Header.Add("ETag", etag.StringValue())
	}

	// TODO: Set timeout handler, and on timeout return 5.04
	resp, err := client.Do(req)
	if err != nil {
		SendMessage(BadGatewayMessage(msg.MessageID, MessageAcknowledgment), session)
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
		SendMessage(BadGatewayMessage(msg.MessageID, MessageAcknowledgment), session)
		return
	}

	_, err = SendMessage(respMsg.GetMessage(), session)
	if err != nil {
		println(err.Error())
	}
}

// Handles requests for proxying from HTTP to CoAP
func HTTPCOAPProxyHandler(msg *Message, conn *net.UDPConn, addr net.Addr) {
	log.Println("HttpCoapProxyHandler Proxy Handler")
}
