package goap

type RouteHandler func(*Message) *Message
type MessageHandler func(*Message)

func (s *Server) NewRoute(path string, method CoapCode, fn RouteHandler) *Route {
	r := &Route{
		AutoAck: false,
		Path:    path,
		Method:  method,
		Handler: fn,
	}
	s.routes = append(s.routes, r)

	return r
}

type Route struct {
	Path       string
	Method     CoapCode
	Handler    RouteHandler
	AutoAck    bool
	MediaTypes []MediaType
}

func (r *Route) AutoAcknowledge(ack bool) *Route {
	r.AutoAck = ack

	return r
}

func (r *Route) BindMediaTypes(ms []MediaType) {
	r.MediaTypes = ms
}
