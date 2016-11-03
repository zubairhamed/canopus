package canopus

import (
	"net"
	"sync"
	"log"
	"github.com/jvermillard/nativedtls"
	"github.com/wendal/errors"
)

func NewClient() *CoapClient {
	return &CoapClient{}
}

type CoapClient struct {

}

func (c *CoapClient) Dial(address string) (conn CoapConnection, err error) {
	udpConn, err := net.Dial("udp", address)
	if err != nil {
		return
	}

	conn = &UDPCoapConnection{
		conn: udpConn,
	}

	return
}

func (c *CoapClient) DialDTLS(address, secret string) (conn CoapConnection, err error) {
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
	dtlsClient.SetPSK("Client_identity", []byte(secret))

	return &DTLSCoapConnection{
		dtlsClient: dtlsClient,
	}, nil

}

type CoapConnection interface {
	ObserveResource(resource string) (tok string, err error)
	CancelObserveResource(resource string, token string) (err error)
	StopObserve(ch chan *ObserveMessage)
	Observe(ch chan *ObserveMessage)
	Send(req CoapRequest) (resp CoapResponse, err error)
	Close()
}

/*
	conn, err := net.Dial("udp", "localhost:5684")
	if err != nil {
		panic(err)
	}
	c := nativedtls.NewDTLSClient(ctx, conn)
	c.SetPSK("Client_identity", []byte("secretPSK"))
	if _, err = c.Write([]byte("Yeah!\n")); err != nil {
		panic(err)
	}

	buff := make([]byte, 1500)
	c.Read(buff)

	fmt.Println("Rcvd:", string(buff))

	c.Close()
 */


type DTLSCoapConnection struct {
	UDPCoapConnection
	dtlsClient *nativedtls.DTLSClient
}

func (c *DTLSCoapConnection) Send(req CoapRequest) (resp CoapResponse, err error) {
	msg := req.GetMessage()
	opt := msg.GetOption(OptionBlock1)

	if opt == nil { // Block1 was not set
		if MessageSizeAllowed(req) != true {
			return nil, ErrMessageSizeTooLongBlockOptionValNotSet
		}
	} else { // Block1 was set
		// log.Println("Block 1 was set")
	}

	if opt != nil {
		blockOpt := Block1OptionFromOption(opt)
		if blockOpt.Value == nil {
			if MessageSizeAllowed(req) != true {
				err = ErrMessageSizeTooLongBlockOptionValNotSet
				return
			} else {
				// - Block # = one and only block (sz = unspecified), whereas 0 = 16bits
				// - MOre bit = 0
			}
		} else {
			payload := msg.Payload.GetBytes()
			payloadLen := uint32(len(payload))
			blockSize := blockOpt.BlockSizeLength()
			currSeq := uint32(0)
			totalBlocks := uint32(payloadLen / blockSize)
			completed := false

			var wg sync.WaitGroup
			wg.Add(1)

			for completed == false {
				if currSeq <= totalBlocks {

					var blockPayloadStart uint32
					var blockPayloadEnd uint32
					var blockPayload []byte

					blockPayloadStart = currSeq*uint32(blockSize) + (currSeq * 1)

					more := true
					if currSeq == totalBlocks {
						more = false
						blockPayloadEnd = payloadLen
					} else {
						blockPayloadEnd = blockPayloadStart + uint32(blockSize)
					}

					blockPayload = payload[blockPayloadStart:blockPayloadEnd]

					blockOpt = NewBlock1Option(blockOpt.Size(), more, currSeq)
					msg.ReplaceOptions(blockOpt.Code, []Option{blockOpt})
					msg.MessageID = GenerateMessageID()
					msg.Payload = NewBytesPayload(blockPayload)

					// send message
					_, err2 := c.sendMessage(msg)
					if err2 != nil {
						wg.Done()
						return
					}
					currSeq = currSeq + 1

				} else {
					completed = true
					wg.Done()
				}
			}
		}
	}

	resp, err = c.sendMessage(msg)
	return
}

func (c *DTLSCoapConnection) sendMessage(msg *Message) (resp CoapResponse, err error) {
	if msg == nil {
		return nil, ErrNilMessage
	}

	b, err := MessageToBytes(msg)
	if err != nil {
		return
	}

	_, err = c.dtlsClient.Write(b)
	if err != nil {
		return
	}

	if msg.MessageType == MessageNonConfirmable {
		resp = NewResponse(NewEmptyMessage(msg.MessageID), nil)
		return
	}

	// c.conn.SetReadDeadline(time.Now().Add(2))

	msgBuf := make([]byte, 1500)
	n, err := c.dtlsClient.Read(msgBuf)
	if err != nil {
		return
	}

	respMsg, err := BytesToMessage(msgBuf[:n])
	if err != nil {
		return
	}
	resp = NewResponse(respMsg, nil)

	return
}

func (c *DTLSCoapConnection) Close() {
	c.dtlsClient.Close()
}


type UDPCoapConnection struct {
	conn net.Conn
}

func (c *UDPCoapConnection) ObserveResource(resource string) (tok string, err error) {
	req := NewRequest(MessageConfirmable, Get, GenerateMessageID())
	req.SetRequestURI(resource)
	req.GetMessage().AddOption(OptionObserve, 0)

	resp, err := c.Send(req)

	tok = string(resp.GetMessage().Token)

	return
}

func (c *UDPCoapConnection) CancelObserveResource(resource string, token string) (err error) {
	req := NewRequest(MessageConfirmable, Get, GenerateMessageID())
	req.SetRequestURI(resource)
	req.GetMessage().Token = []byte(token)
	req.GetMessage().AddOption(OptionObserve, 1)
	_, err = c.Send(req)
	return
}

func (c *UDPCoapConnection) StopObserve(ch chan *ObserveMessage) {
	log.Println("StopObserve")
	close(ch)
}

func (c *UDPCoapConnection) Close() {
	c.conn.Close()
}

func (c *UDPCoapConnection) Observe(ch chan *ObserveMessage)  {
	conn := c.conn

	readBuf := make([]byte, MaxPacketSize)
	for {
		len, err := conn.Read(readBuf)
		if err == nil {
			msgBuf := make([]byte, len)
			copy(msgBuf, readBuf)

			msg, err := BytesToMessage(msgBuf)
			if msg.GetOption(OptionObserve) != nil {
				ch <- NewObserveMessage(msg.GetURIPath(), msg.Payload, msg)
			}
			if err != nil {
				log.Println("Error occured reading UDP", err)
				close(ch)
			}
		} else {
			log.Println("Error occured reading UDP", err)
			close(ch)
		}
	}
}

func (c *UDPCoapConnection) Send(req CoapRequest) (resp CoapResponse, err error) {
	msg := req.GetMessage()
	opt := msg.GetOption(OptionBlock1)

	if opt == nil { // Block1 was not set
		if MessageSizeAllowed(req) != true {
			return nil, ErrMessageSizeTooLongBlockOptionValNotSet
		}
	} else { // Block1 was set
		// log.Println("Block 1 was set")
	}

	if opt != nil {
		blockOpt := Block1OptionFromOption(opt)
		if blockOpt.Value == nil {
			if MessageSizeAllowed(req) != true {
				err = ErrMessageSizeTooLongBlockOptionValNotSet
				return
			} else {
				// - Block # = one and only block (sz = unspecified), whereas 0 = 16bits
				// - MOre bit = 0
			}
		} else {
			payload := msg.Payload.GetBytes()
			payloadLen := uint32(len(payload))
			blockSize := blockOpt.BlockSizeLength()
			currSeq := uint32(0)
			totalBlocks := uint32(payloadLen / blockSize)
			completed := false

			var wg sync.WaitGroup
			wg.Add(1)

			for completed == false {
				if currSeq <= totalBlocks {

					var blockPayloadStart uint32
					var blockPayloadEnd uint32
					var blockPayload []byte

					blockPayloadStart = currSeq*uint32(blockSize) + (currSeq * 1)

					more := true
					if currSeq == totalBlocks {
						more = false
						blockPayloadEnd = payloadLen
					} else {
						blockPayloadEnd = blockPayloadStart + uint32(blockSize)
					}

					blockPayload = payload[blockPayloadStart:blockPayloadEnd]

					blockOpt = NewBlock1Option(blockOpt.Size(), more, currSeq)
					msg.ReplaceOptions(blockOpt.Code, []Option{blockOpt})
					msg.MessageID = GenerateMessageID()
					msg.Payload = NewBytesPayload(blockPayload)

					// send message
					_, err2 := c.sendMessage(msg)
					if err2 != nil {
						wg.Done()
						return
					}
					currSeq = currSeq + 1

				} else {
					completed = true
					wg.Done()
				}
			}
		}
	}

	resp, err = c.sendMessage(msg)
	return
}

func (c *UDPCoapConnection) sendMessage(msg *Message) (resp CoapResponse, err error) {
	if msg == nil {
		return nil, ErrNilMessage
	}

	b, err := MessageToBytes(msg)
	if err != nil {
		return
	}

	_, err = c.conn.Write(b)
	if err != nil {
		return
	}

	if msg.MessageType == MessageNonConfirmable {
		resp = NewResponse(NewEmptyMessage(msg.MessageID), nil)
		return
	}

	// c.conn.SetReadDeadline(time.Now().Add(2))

	msgBuf := make([]byte, 1500)
	n, err := c.conn.Read(msgBuf)
	if err != nil {
		return
	}

	respMsg, err := BytesToMessage(msgBuf[:n])
	if err != nil {
		return
	}
	resp = NewResponse(respMsg, nil)

	return
}

func NewObserveMessage(r string, val interface{}, msg *Message) *ObserveMessage {
	return &ObserveMessage{
		Resource: r,
		Value: val,
		Msg: msg,
	}
}

type ObserveMessage struct {
	Resource string
	Value interface{}
	Msg *Message
}
