package canopus

import (
	"net"

	"github.com/jvermillard/nativedtls"
	"github.com/wendal/errors"
)

func NewClient() Client {
	return &CoapClient{}
}

type CoapClient struct {
}

func (c *CoapClient) Dial(address string) (conn ClientConnection, err error) {
	udpConn, err := net.Dial("udp", address)
	if err != nil {
		return
	}

	conn = &UDPClientConnection{
		conn: udpConn,
	}

	return
}

func (c *CoapClient) DialDTLS(address, psk string) (conn ClientConnection, err error) {
	ctx := nativedtls.NewDTLSContext()
	if !ctx.SetCipherList("PSK-AES256-CCM8:PSK-AES128-CCM8") {
		err = errors.New("impossible to set cipherlist")
		return
	}

	udpConn, err := net.Dial("udp", address)
	if err != nil {
		return
	}

	dtlsClient := nativedtls.NewDTLSClient(ctx, udpConn)
	dtlsClient.SetPSK("Client_identity", []byte(psk))

	conn = &DTLSClientConnection{
		dtlsClient: dtlsClient,
	}

	return
}

func NewObserveMessage(r string, val interface{}, msg Message) ObserveMessage {
	return &CoapObserveMessage{
		Resource: r,
		Value:    val,
		Msg:      msg,
	}
}

type CoapObserveMessage struct {
	CoapMessage
	Resource string
	Value    interface{}
	Msg      Message
}

func (m *CoapObserveMessage) GetResource() string {
	return m.Resource
}

func (m *CoapObserveMessage) GetValue() interface{} {
	return m.Value
}

func (m *CoapObserveMessage) GetMessage() Message {
	return m.GetMessage()
}
