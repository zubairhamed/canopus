package canopus

// Represents a message payload containing an array of bytes
func NewBytesPayload(v []byte) MessagePayload {
	return &BytesPayload{
		content: v,
	}
}

type BytesPayload struct {
	content []byte
}

func (p *BytesPayload) GetBytes() []byte {
	return p.content
}

func (p *BytesPayload) Length() int {
	return len(p.content)
}

func (p *BytesPayload) String() string {
	return string(p.content)
}
