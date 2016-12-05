package canopus

import (
	"encoding/json"
	"log"
)

func NewJSONPayload(obj interface{}) MessagePayload {
	return &JSONPayload{
		obj: obj,
	}
}

// Represents a message payload containing JSON String
type JSONPayload struct {
	obj interface{}
}

func (p *JSONPayload) GetBytes() []byte {
	o, err := json.MarshalIndent(p.obj, "", "   ")

	if err != nil {
		log.Println(err)

		return []byte{}
	}

	return []byte(string(o))
}

func (p *JSONPayload) Length() int {
	return 0
}

func (p *JSONPayload) String() string {
	o, _ := json.Marshal(p.obj)

	return string(o)
}
