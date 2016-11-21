package canopus

import (
	"fmt"
	"log"
	"net"
	"sync"
)

func MessageSizeAllowed(req Request) bool {
	msg := req.GetMessage()
	b, _ := MessageToBytes(msg)

	if len(b) > 65536 {
		return false
	}

	return true
}

type UDPConnection struct {
	conn net.Conn
}

func (c *UDPConnection) ObserveResource(resource string) (tok string, err error) {
	req := NewRequest(MessageConfirmable, Get, GenerateMessageID())
	req.SetRequestURI(resource)
	req.GetMessage().AddOption(OptionObserve, 0)

	resp, err := c.Send(req)
	tok = string(resp.GetMessage().GetToken())

	return
}

func (c *UDPConnection) CancelObserveResource(resource string, token string) (err error) {
	req := NewRequest(MessageConfirmable, Get, GenerateMessageID())
	req.SetRequestURI(resource)
	req.GetMessage().AddOption(OptionObserve, 1)

	_, err = c.Send(req)
	return
}

func (c *UDPConnection) StopObserve(ch chan ObserveMessage) {
	close(ch)
}

func (c *UDPConnection) Close() error {
	return c.conn.Close()
}

func (c *UDPConnection) Observe(ch chan ObserveMessage) {

	readBuf := make([]byte, MaxPacketSize)
	for {
		len, err := c.Read(readBuf)
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

func (c *UDPConnection) Send(req Request) (resp Response, err error) {
	fmt.Println("UDPConnection:SendA")
	msg := req.GetMessage()
	opt := msg.GetOption(OptionBlock1)

	fmt.Println("UDPConnection:SendB")
	if opt == nil { // Block1 was not set
		fmt.Println("UDPConnection:SendC")
		if MessageSizeAllowed(req) != true {
			return nil, ErrMessageSizeTooLongBlockOptionValNotSet
		}
	} else { // Block1 was set
		// log.Println("Block 1 was set")
	}

	fmt.Println("UDPConnection:SendD")
	if opt != nil {
		fmt.Println("UDPConnection:SendE")
		blockOpt := Block1OptionFromOption(opt)
		fmt.Println("UDPConnection:SendF")
		if blockOpt.Value == nil {
			fmt.Println("UDPConnection:SendG")
			if MessageSizeAllowed(req) != true {
				err = ErrMessageSizeTooLongBlockOptionValNotSet
				return
			} else {
				fmt.Println("UDPConnection:SendH")
				// - Block # = one and only block (sz = unspecified), whereas 0 = 16bits
				// - MOre bit = 0
			}
		} else { // BLock transfer request
			fmt.Println("UDPConnection:SendI")
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
					modifiedMsg := msg.(*CoapMessage)
					modifiedMsg.SetMessageId(GenerateMessageID())
					modifiedMsg.SetPayload(NewBytesPayload(blockPayload))

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
	fmt.Println("UDPConnection:SendJ")
	resp, err = c.sendMessage(msg)
	fmt.Println("UDPConnection:SendK")
	return
}

func (c *UDPConnection) sendMessage(msg Message) (resp Response, err error) {
	fmt.Println("UDPConnection:sendMessageA")
	if msg == nil {
		fmt.Println("UDPConnection:sendMessageB")
		return nil, ErrNilMessage
	}
	fmt.Println("UDPConnection:sendMessageC")

	fmt.Println("UDPConnection:sendMessageD")
	b, err := MessageToBytes(msg)
	fmt.Println("UDPConnection:sendMessageE")
	if err != nil {
		fmt.Println("UDPConnection:sendMessageF")
		return
	}
	fmt.Println("UDPConnection:sendMessageG")

	_, err = c.Write(b)
	fmt.Println("UDPConnection:sendMessageH")
	if err != nil {
		fmt.Println("UDPConnection:sendMessageI")
		return
	}

	fmt.Println("UDPConnection:sendMessageK")
	if msg.GetMessageType() == MessageNonConfirmable {
		fmt.Println("UDPConnection:sendMessageL")
		resp = NewResponse(NewEmptyMessage(msg.GetMessageId()), nil)
		return
	}

	fmt.Println("UDPConnection:sendMessageM")
	// c.conn.SetReadDeadline(time.Now().Add(2))

	msgBuf := make([]byte, 1500)
	fmt.Println("UDPConnection:sendMessageN")
	if msg.GetMessageType() == MessageAcknowledgment {
		go c.Read(msgBuf)
		return
	}
	n, err := c.Read(msgBuf)
	fmt.Println("UDPConnection:sendMessageO")
	if err != nil {
		return
	}

	fmt.Println("UDPConnection:sendMessageP")
	respMsg, err := BytesToMessage(msgBuf[:n])
	fmt.Println("UDPConnection:sendMessageQ")
	if err != nil {
		return
	}
	fmt.Println("UDPConnection:sendMessageR")
	resp = NewResponse(respMsg, nil)

	fmt.Println("UDPConnection:sendMessageS")
	if msg.GetMessageType() == MessageConfirmable {
		fmt.Println("UDPConnection:sendMessageT")
		ack := NewMessageOfType(MessageAcknowledgment, msg.GetMessageId(), nil)

		fmt.Println("UDPConnection:sendMessageU")
		c.Send(NewRequestFromMessage(ack))
	}
	return
}

func (c *UDPConnection) Write(b []byte) (int, error) {
	fmt.Println("UDPConnection:Write")
	return c.conn.Write(b)
}

func (c *UDPConnection) Read(b []byte) (int, error) {
	return c.conn.Read(b)
}
