package canopus

func NewEvent(data map[string]interface{}) *Event {
	return &Event{
		Data: data,
	}
}

type Event struct {
	Message *Message
	Data    map[string]interface{}
}
