package canopus

import (
	"bytes"
	"encoding/json"
	"log"
)

type MessagePayload interface {
	GetBytes() []byte
	Length() int
	String() string
}

func NewPlainTextPayload(s string) MessagePayload {
	return &PlainTextPayload{
		content: s,
	}
}

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

type XmlPayload struct {
}

func NewEmptyPayload() MessagePayload {
	return &EmptyPayload{}
}

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
