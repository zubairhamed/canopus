package canopus

func NewEmptyPayload() MessagePayload {
	return &EmptyPayload{}
}

// Represents an empty message payload
type EmptyPayload struct {
}

func (p *EmptyPayload) GetBytes() []byte {
	return []byte{}
}

func (p *EmptyPayload) Length() int {
	return 0
}

func (p *EmptyPayload) String() string {
	return ""
}
