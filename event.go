package goap

type EventCode uint8

const (
	EVENT_SERVER_STARTED = EventCode(0)
	EVENT_DISCOVERY      = EventCode(1)
)

func NewEvent(data map[string]interface{}) *Event {
	return &Event{
		Data: data,
	}
}

type Event struct {
	Message *Message
	Data    map[string]interface{}
}
