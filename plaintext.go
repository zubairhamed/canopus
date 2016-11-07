package canopus

import "bytes"

// Instantiates a new message payload of type string
func NewPlainTextPayload(s string) MessagePayload {
	return &PlainTextPayload{
		content: s,
	}
}

// Represents a message payload containing string value
type PlainTextPayload struct {
	content string
}

func (p *PlainTextPayload) GetBytes() []byte {
	return bytes.NewBufferString(p.content).Bytes()
}

func (p *PlainTextPayload) Length() int {
	return len(p.content)
}

func (p *PlainTextPayload) String() string {
	return p.content
}
