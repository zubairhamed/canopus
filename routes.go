package goap

type RouteHandler func(*Message) *Message

type Route struct {
	Path		string
	methods 	[]uint8
	Handler 	RouteHandler
	AutoAck 	bool
}

func (r *Route) Methods(methods []string) {

}

func (r *Route) AutoAcknowledge(ack bool) {
	r.AutoAck = ack
}
