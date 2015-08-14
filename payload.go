package canopus

import (
	"bytes"
	"encoding/json"
	"log"
)

// Represents the payload/content of a CoAP Message
type MessagePayload interface {
	GetBytes() []byte
	Length() int
	String() string
}

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

// Represents a message payload containing XML String
type XmlPayload struct {
}

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

func NewJsonPayload(obj interface{}) MessagePayload {
	return &JsonPayload{
		obj: obj,
	}
}

// Represents a message payload containing JSON String
type JsonPayload struct {
	obj interface{}
}

func (p *JsonPayload) GetBytes() []byte {
	o, err := json.MarshalIndent(p.obj, "", "   ")

	if err != nil {
		log.Println(err)

		return []byte{}
	}

	return []byte(string(o))
}

func (p *JsonPayload) Length() int {
	return 0
}

func (p *JsonPayload) String() string {
	o, _ := json.Marshal(p.obj)

	return string(o)
}
