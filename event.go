package canopus

type FnEventNotify func(string, interface{}, *Message)
type FnEventStart func(CoapServer)
type FnEventClose func(CoapServer)
type FnEventDiscover func()
type FnEventError func(error)
type FnEventObserve func(string, *Message)
type FnEventObserveCancel func(string, *Message)
type FnEventMessage func(*Message, bool)

type EventCode int

const (
	EventStart         EventCode = 0
	EventClose         EventCode = 1
	EventDiscover      EventCode = 2
	EventMessage       EventCode = 3
	EventError         EventCode = 4
	EventObserve       EventCode = 5
	EventObserveCancel EventCode = 6
	EventNotify        EventCode = 7
)

func NewEvents() *Events {
	return &Events{
		evtFnNotify:        []FnEventNotify{},
		evtFnStart:         []FnEventStart{},
		evtFnClose:         []FnEventClose{},
		evtFnDiscover:      []FnEventDiscover{},
		evtFnError:         []FnEventError{},
		evtFnObserve:       []FnEventObserve{},
		evtFnObserveCancel: []FnEventObserveCancel{},
		evtFnMessage:       []FnEventMessage{},
	}
}

// This holds the various events which are triggered throughout
// an application's lifetime
type Events struct {
	evtFnNotify        []FnEventNotify
	evtFnStart         []FnEventStart
	evtFnClose         []FnEventClose
	evtFnDiscover      []FnEventDiscover
	evtFnError         []FnEventError
	evtFnObserve       []FnEventObserve
	evtFnObserveCancel []FnEventObserveCancel
	evtFnMessage       []FnEventMessage
}

// OnNotify is Fired when an observeed resource is notified
func (ce *Events) OnNotify(fn FnEventNotify) {
	ce.evtFnNotify = append(ce.evtFnNotify, fn)
}

// Fired when the server/client starts up
func (ce *Events) OnStart(fn FnEventStart) {
	ce.evtFnStart = append(ce.evtFnStart, fn)
}

// Fired when the server/client closes
func (ce *Events) OnClose(fn FnEventClose) {
	ce.evtFnClose = append(ce.evtFnClose, fn)
}

// Fired when a discovery request is triggered
func (ce *Events) OnDiscover(fn FnEventDiscover) {
	ce.evtFnDiscover = append(ce.evtFnDiscover, fn)
}

// Catch-all event which is fired when an error occurs
func (ce *Events) OnError(fn FnEventError) {
	ce.evtFnError = append(ce.evtFnError, fn)
}

// Fired when an observe request is triggered for a resource
func (ce *Events) OnObserve(fn FnEventObserve) {
	ce.evtFnObserve = append(ce.evtFnObserve, fn)
}

// Fired when an observe-cancel request is triggered foa r esource
func (ce *Events) OnObserveCancel(fn FnEventObserveCancel) {
	ce.evtFnObserveCancel = append(ce.evtFnObserveCancel, fn)
}

// Fired when a message is received
func (ce *Events) OnMessage(fn FnEventMessage) {
	ce.evtFnMessage = append(ce.evtFnMessage, fn)
}

// Fires the "OnNotify" event
func (ce *Events) Notify(resource string, value interface{}, msg *Message) {
	for _, fn := range ce.evtFnNotify {
		fn(resource, value, msg)
	}
}

// Fires the "OnStarted" event
func (ce *Events) Started(server CoapServer) {
	for _, fn := range ce.evtFnStart {
		fn(server)
	}
}

// Fires the "OnClosed" event
func (ce *Events) Closed(server CoapServer) {
	for _, fn := range ce.evtFnClose {
		fn(server)
	}
}

// Fires the "OnDiscover" event
func (ce *Events) Discover() {
	for _, fn := range ce.evtFnDiscover {
		fn()
	}
}

// Fires the "OnError" event given an error object
func (ce *Events) Error(err error) {
	for _, fn := range ce.evtFnError {
		fn(err)
	}
}

// Fires the "OnObserve" event for a given resource
func (ce *Events) Observe(resource string, msg *Message) {
	for _, fn := range ce.evtFnObserve {
		fn(resource, msg)
	}
}

// Fires the "OnObserveCancelled" event for a given resource
func (ce *Events) ObserveCancelled(resource string, msg *Message) {
	for _, fn := range ce.evtFnObserveCancel {
		fn(resource, msg)
	}
}

// Fires the "OnMessage" event. The 'inbound' variables determines if the
// message is inbound or outgoing
func (ce *Events) Message(msg *Message, inbound bool) {
	for _, fn := range ce.evtFnMessage {
		fn(msg, inbound)
	}
}
