package goap

type RouteHandler func(*Message) *Message


func (s *Server) NewRoute(path string, method uint8, fn RouteHandler,) (*Route) {
	r := &Route{
		AutoAck: true,
		Path: path,
		Method: method,
		Handler: fn,
	}
	s.routes = append(s.routes, r)

	return r
}

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
