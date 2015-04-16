package goap


type EventCode uint8
const (
    EVENT_SERVER_STARTED = EventCode(0)
    EVENT_DISCOVERY = EventCode(1)
)

func NewEvent() (*Event) {
    return &Event{}
}

type Event struct {
    Message     *Message
}
