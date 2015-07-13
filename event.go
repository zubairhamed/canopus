package canopus

type FnEventNotify func(string, interface{}, *Message)
type FnEventStart func(*CoapServer)
type FnEventClose func(*CoapServer)
type FnEventDiscover func()
type FnEventError func(error)
type FnEventObserve func(string, *Message)
type FnEventObserveCancel func(string, *Message)
type FnEventMessage func(*Message, bool)

type EventCode int
const (
	EVT_START          EventCode = 0
	EVT_CLOSE          EventCode = 1
	EVT_DISCOVER       EventCode = 2
	EVT_MESSAGE        EventCode = 3
	EVT_ERROR          EventCode = 4
	EVT_OBSERVE        EventCode = 5
	EVT_OBSERVE_CANCEL EventCode = 6
	EVT_NOTIFY		   EventCode = 7
)

func NewCanopusEvents() *CanopusEvents {
	return &CanopusEvents{
		evtFnNotify: 		[]FnEventNotify{},
		evtFnStart: 		[]FnEventStart{},
		evtFnClose: 		[]FnEventClose{},
		evtFnDiscover: 		[]FnEventDiscover{},
		evtFnError: 		[]FnEventError{},
		evtFnObserve: 		[]FnEventObserve{},
		evtFnObserveCancel: []FnEventObserveCancel{},
		evtFnMessage: 		[]FnEventMessage{},
	}
}
type CanopusEvents struct {
	evtFnNotify 		[]FnEventNotify
	evtFnStart			[]FnEventStart
	evtFnClose			[]FnEventClose
	evtFnDiscover		[]FnEventDiscover
	evtFnError			[]FnEventError
	evtFnObserve		[]FnEventObserve
	evtFnObserveCancel	[]FnEventObserveCancel
	evtFnMessage		[]FnEventMessage
}

func (ce *CanopusEvents) OnNotify(fn FnEventNotify) {
	ce.evtFnNotify = append(ce.evtFnNotify, fn)
}

func (ce *CanopusEvents) OnStart(fn FnEventStart) {
	ce.evtFnStart = append(ce.evtFnStart, fn)
}

func (ce *CanopusEvents) OnClose(fn FnEventClose) {
	ce.evtFnClose = append(ce.evtFnClose, fn)
}

func (ce *CanopusEvents) OnDiscover(fn FnEventDiscover) {
	ce.evtFnDiscover = append(ce.evtFnDiscover, fn)
}

func (ce *CanopusEvents) OnError(fn FnEventError) {
	ce.evtFnError = append(ce.evtFnError, fn)
}

func (ce *CanopusEvents) OnObserve(fn FnEventObserve) {
	ce.evtFnObserve = append(ce.evtFnObserve, fn)
}

func (ce *CanopusEvents) OnObserveCancel(fn FnEventObserveCancel){
	ce.evtFnObserveCancel = append(ce.evtFnObserveCancel, fn)
}

func (ce *CanopusEvents) OnMessage(fn FnEventMessage){
	ce.evtFnMessage = append(ce.evtFnMessage, fn)
}

func (ce *CanopusEvents) Notify(resource string, value interface{}, msg *Message){
	for _, fn := range ce.evtFnNotify {
		fn(resource, value, msg)
	}
}

func (ce *CanopusEvents) Started(server *CoapServer){
	for _, fn := range ce.evtFnStart {
		fn(server)
	}
}

func (ce *CanopusEvents) Closed(server *CoapServer){
	for _, fn := range ce.evtFnClose {
		fn(server)
	}
}

func (ce *CanopusEvents) Discover(){
	for _, fn := range ce.evtFnDiscover {
		fn()
	}
}

func (ce *CanopusEvents) Error(err error){
	for _, fn := range ce.evtFnError {
		fn(err)
	}
}

func (ce *CanopusEvents) Observe(resource string, msg *Message){
	for _, fn := range ce.evtFnObserve {
		fn(resource, msg)
	}
}

func (ce *CanopusEvents) ObserveCancelled(resource string, msg *Message){
	for _, fn := range ce.evtFnObserveCancel {
		fn(resource, msg)
	}
}

func (ce *CanopusEvents) Message(msg *Message, inbound bool){
	for _, fn := range ce.evtFnMessage {
		fn(msg, inbound)
	}
}


