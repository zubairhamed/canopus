package canopus

// Represents a message payload containing core-link format values
type CoreLinkFormatPayload struct {
}

func (p *CoreLinkFormatPayload) GetBytes() []byte {
	return make([]byte, 0)
}

func (p *CoreLinkFormatPayload) Length() int {
	return 0
}

func (p *CoreLinkFormatPayload) String() string {
	return ""
}
