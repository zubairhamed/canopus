package canopus

import (
	"log"
	"net"
	"sync"

	"github.com/jvermillard/nativedtls"
)

type DTLSClientConnection struct {
	UDPClientConnection
	dtlsClient *nativedtls.DTLSClient
}

func (c *DTLSClientConnection) Send(req Request) (resp Response, err error) {
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
			payload := msg.GetPayload().GetBytes()
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
					msg.SetMessageId(GenerateMessageID())
					msg.SetPayload(NewBytesPayload(blockPayload))

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

func (c *DTLSClientConnection) sendMessage(msg Message) (resp Response, err error) {
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

	if msg.GetMessageType() == MessageNonConfirmable {
		resp = NewResponse(NewEmptyMessage(msg.GetMessageId()), nil)
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

func (c *DTLSClientConnection) Close() {
	c.dtlsClient.Close()
}

func MessageSizeAllowed(req Request) bool {
	msg := req.GetMessage()
	b, _ := MessageToBytes(msg)

	if len(b) > 65536 {
		return false
	}

	return true
}

type UDPClientConnection struct {
	conn net.Conn
}

func (c *UDPClientConnection) ObserveResource(resource string) (tok string, err error) {
	req := NewRequest(MessageConfirmable, Get, GenerateMessageID())
	req.SetRequestURI(resource)
	req.GetMessage().AddOption(OptionObserve, 0)

	resp, err := c.Send(req)
	tok = string(resp.GetMessage())

	return
}

func (c *UDPClientConnection) CancelObserveResource(resource string, token string) (err error) {
	req := NewRequest(MessageConfirmable, Get, GenerateMessageID())
	req.SetRequestURI(resource)
	req.GetMessage().AddOption(OptionObserve, 0)

	_, err = c.Send(req)
	return
}

func (c *UDPClientConnection) StopObserve(ch chan *ObserveMessage) {
	close(ch)
}

func (c *UDPClientConnection) Close() {
	c.conn.Close()
}

func (c *UDPClientConnection) Observe(ch chan *ObserveMessage) {
	conn := c.conn

	readBuf := make([]byte, MaxPacketSize)
	for {
		len, err := conn.Read(readBuf)
		if err == nil {
			msgBuf := make([]byte, len)
			copy(msgBuf, readBuf)

			msg, err := BytesToMessage(msgBuf)
			if msg.GetOption(OptionObserve) != nil {
				ch <- NewObserveMessage(msg.GetURIPath(), msg.GetPayload(), msg)
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

func (c *UDPClientConnection) Send(req Request) (resp Response, err error) {
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
		} else { // BLock transfer request
			payload := msg.GetPayload().GetBytes()
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
					msg.SetMessageId(GenerateMessageID())
					msg.SetPayload(NewBytesPayload(blockPayload))

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

func (c *UDPClientConnection) sendMessage(msg Message) (resp Response, err error) {
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

	if msg.GetMessageType() == MessageNonConfirmable {
		resp = NewResponse(NewEmptyMessage(msg.GetMessageId()), nil)
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

	if msg.GetMessageType() == MessageConfirmable {
		ack := NewMessageOfType(MessageAcknowledgment, msg.GetMessageId())

		c.Send(NewRequestFromMessage(ack))
	}

	return
}
