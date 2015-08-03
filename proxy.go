package canopus

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
)

type ProxyHandler func(msg *Message, conn *net.UDPConn, addr *net.UDPAddr)

func NullProxyHandler(msg *Message, conn *net.UDPConn, addr *net.UDPAddr) {
	SendMessageTo(NotProxyingSupportedMessage(msg.MessageId), conn, addr)
}

func CoapCoapProxyHandler(msg *Message, conn *net.UDPConn, addr *net.UDPAddr) {
	log.Println("CoapCoapProxyHandler Proxy Handler")
}

func CoapHttpProxyHandler(msg *Message, conn *net.UDPConn, addr *net.UDPAddr) {
	proxyUri := msg.GetOption(OPTION_PROXY_URI).StringValue()

	client := &http.Client{}

	req, _ := http.NewRequest(MethodString(CoapCode(msg.GetMethod())), proxyUri, nil)

	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		log.Println(err)
	}

	contents, _ := ioutil.ReadAll(resp.Body)
	msg.Payload = NewBytesPayload(contents)

	respMsg := NewRequestFromMessage(msg)
	SendMessageTo(respMsg.GetMessage(), conn, addr)

	log.Println("CoapHttpProxyHandler Proxy Handler, Call:", proxyUri)
}

func HttpCoapProxyHandler(msg *Message, conn *net.UDPConn, addr *net.UDPAddr) {
	log.Println("HttpCoapProxyHandler Proxy Handler")
}
