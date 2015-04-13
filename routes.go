package goap
import (
    "regexp"
    "strings"
)

func (s *Server) NewRoute(path string, method CoapCode, fn RouteHandler) *Route {
    re, _ := regexp.Compile(`{[a-z]+}`)
    matches := re.FindAllStringSubmatch(path, -1)

    path = "^" + path
    for _, b := range matches {
        origAttr := b[0]
        attr := strings.Replace(strings.Replace(origAttr, "{", "", -1), "}", "", -1)
        frag := `(?P<` + attr + `>\w+)`
        path = strings.Replace(path, origAttr, frag, -1)
    }
    path += "$"
    re, _ = regexp.Compile(path)

    /*
    OnNewRoute
        Get all values between #{ }
        Construct New RegEx
            Create SubGroups
            Escape any RegEx Values
        Compile and Store Compiled RegEx

    */

	r := &Route{
		AutoAck: false,
		Path:    path,
		Method:  method,
		Handler: fn,
        RegEx: re,
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
    RegEx      *regexp.Regexp
}

func (r *Route) AutoAcknowledge(ack bool) *Route {
	r.AutoAck = ack

	return r
}

func (r *Route) BindMediaTypes(ms []MediaType) {
	r.MediaTypes = ms
}

func (r *Route) Matches(s string) (bool, map[string]string) {
    matches := r.RegEx.FindAllStringSubmatch(s, -1)
    attrs := make(map[string]string)
    if len(matches) > 0 {
        subExp := r.RegEx.SubexpNames()
        for idx, exp := range subExp {
            attrs[exp] = matches[0][idx]
        }

        return true, attrs
    }

    return false, attrs
}
