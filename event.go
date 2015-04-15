package goap


type EventCode uint8
const (
    EVENT_DISCOVERY = EventCode(0)
)

func NewEvent() (*Event) {
    return &Event{}
}

type Event struct {
    Message     *Message
}

