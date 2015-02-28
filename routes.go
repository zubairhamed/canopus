package goap

type RouteHandler func(*Message) *Message

type Route struct {
	Path		string
	Method 		uint8
	Handler 	RouteHandler
	AutoAck 	bool
}

func (r *Route) AutoAcknowledge(ack bool) (*Route) {
	r.AutoAck = ack

	return r
}
