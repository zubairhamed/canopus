package goap

import "bytes"

type MessagePayload interface {
    GetBytes() ([]byte)
    Length() (int)
    ToString() (string)
}

func NewPlainTextPayload(s string) (MessagePayload) {
    return &PlainTextPayload{
        content: s,
    }
}

type PlainTextPayload struct {
    content   string
}

func (p *PlainTextPayload) GetBytes() ([]byte) {
    return bytes.NewBufferString(p.content).Bytes()
}

func (p *PlainTextPayload) Length() (int) {
    return len(p.content)
}

func (p *PlainTextPayload) ToString() (string) {
    return p.content
}

type CoreLinkFormatPayload struct {

}

func (p *CoreLinkFormatPayload) GetBytes() ([]byte) {
    return make([]byte, 0)
}

func (p *CoreLinkFormatPayload) Length() (int) {
    return 0
}

func (p *CoreLinkFormatPayload) ToString() (string) {
    return ""
}

func NewBytesPayload(b []byte) (MessagePayload) {
    return &BytesPayload{
        content: b,
    }
}

type BytesPayload struct {
    content     []byte
}

func (p *BytesPayload) GetBytes() ([]byte) {
    return make([]byte, 0)
}

func (p *BytesPayload) Length() (int) {
    return len(p.content)
}

func (p *BytesPayload) ToString() (string) {
    return string(p.content)
}

type XmlPayload struct {

}

type EmptyPayload struct {

}

