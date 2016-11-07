package canopus

func NewEvents() *ServerEvents {
	return &ServerEvents{
		evtFnNotify:        []FnEventNotify{},
		evtFnStart:         []FnEventStart{},
		evtFnClose:         []FnEventClose{},
		evtFnDiscover:      []FnEventDiscover{},
		evtFnError:         []FnEventError{},
		evtFnObserve:       []FnEventObserve{},
		evtFnObserveCancel: []FnEventObserveCancel{},
		evtFnMessage:       []FnEventMessage{},
		evtFnBlockMessage:  []FnEventBlockMessage{},
	}
}

// This holds the various events which are triggered throughout
// an application's lifetime
type ServerEvents struct {
	evtFnNotify        []FnEventNotify
	evtFnStart         []FnEventStart
	evtFnClose         []FnEventClose
	evtFnDiscover      []FnEventDiscover
	evtFnError         []FnEventError
	evtFnObserve       []FnEventObserve
	evtFnObserveCancel []FnEventObserveCancel
	evtFnMessage       []FnEventMessage
	evtFnBlockMessage  []FnEventBlockMessage
}

// OnNotify is Fired when an observeed resource is notified
func (ce *ServerEvents) OnNotify(fn FnEventNotify) {
	ce.evtFnNotify = append(ce.evtFnNotify, fn)
}

// Fired when the server/client starts up
func (ce *ServerEvents) OnStart(fn FnEventStart) {
	ce.evtFnStart = append(ce.evtFnStart, fn)
}

// Fired when the server/client closes
func (ce *ServerEvents) OnClose(fn FnEventClose) {
	ce.evtFnClose = append(ce.evtFnClose, fn)
}

// Fired when a discovery request is triggered
func (ce *ServerEvents) OnDiscover(fn FnEventDiscover) {
	ce.evtFnDiscover = append(ce.evtFnDiscover, fn)
}

// Catch-all event which is fired when an error occurs
func (ce *ServerEvents) OnError(fn FnEventError) {
	ce.evtFnError = append(ce.evtFnError, fn)
}

// Fired when an observe request is triggered for a resource
func (ce *ServerEvents) OnObserve(fn FnEventObserve) {
	ce.evtFnObserve = append(ce.evtFnObserve, fn)
}

// Fired when an observe-cancel request is triggered foa r esource
func (ce *ServerEvents) OnObserveCancel(fn FnEventObserveCancel) {
	ce.evtFnObserveCancel = append(ce.evtFnObserveCancel, fn)
}

// Fired when a message is received
func (ce *ServerEvents) OnMessage(fn FnEventMessage) {
	ce.evtFnMessage = append(ce.evtFnMessage, fn)
}

// Fired when a block messageis received
func (ce *ServerEvents) OnBlockMessage(fn FnEventBlockMessage) {
	ce.evtFnBlockMessage = append(ce.evtFnBlockMessage, fn)
}

// Fires the "OnNotify" event
func (ce *ServerEvents) Notify(resource string, value interface{}, msg Message) {
	for _, fn := range ce.evtFnNotify {
		fn(resource, value, msg)
	}
}

// Fires the "OnStarted" event
func (ce *ServerEvents) Started(server CoapServer) {
	for _, fn := range ce.evtFnStart {
		fn(server)
	}
}

// Fires the "OnClosed" event
func (ce *ServerEvents) Closed(server CoapServer) {
	for _, fn := range ce.evtFnClose {
		fn(server)
	}
}

// Fires the "OnDiscover" event
func (ce *ServerEvents) Discover() {
	for _, fn := range ce.evtFnDiscover {
		fn()
	}
}

// Fires the "OnError" event given an error object
func (ce *ServerEvents) Error(err error) {
	for _, fn := range ce.evtFnError {
		fn(err)
	}
}

// Fires the "OnObserve" event for a given resource
func (ce *ServerEvents) Observe(resource string, msg Message) {
	for _, fn := range ce.evtFnObserve {
		fn(resource, msg)
	}
}

// Fires the "OnObserveCancelled" event for a given resource
func (ce *ServerEvents) ObserveCancelled(resource string, msg Message) {
	for _, fn := range ce.evtFnObserveCancel {
		fn(resource, msg)
	}
}

// Fires the "OnMessage" event. The 'inbound' variables determines if the
// message is inbound or outgoing
func (ce *ServerEvents) Message(msg Message, inbound bool) {
	for _, fn := range ce.evtFnMessage {
		fn(msg, inbound)
	}
}

// Fires the "OnBlockMessage" event. The 'inbound' variables determines if the
// message is inbound or outgoing
func (ce *ServerEvents) BlockMessage(msg Message, inbound bool) {
	for _, fn := range ce.evtFnBlockMessage {
		fn(msg, inbound)
	}
}
