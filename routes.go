package canopus

import (
	"fmt"
	"regexp"
)

// CreateNewRoute creates a new Route object
func CreateNewRegExRoute(path string, method string, fn RouteHandler) Route {
	var re *regexp.Regexp
	regexpString := path

	// Dots
	re = regexp.MustCompile(`([^\\])\.`)
	regexpString = re.ReplaceAllStringFunc(regexpString, func(m string) string {
		return fmt.Sprintf(`%s\.`, string(m[0]))
	})

	// Wildcard names
	re = regexp.MustCompile(`:[^/#?()\.\\]+\*`)
	regexpString = re.ReplaceAllStringFunc(regexpString, func(m string) string {
		return fmt.Sprintf("(?P<%s>.+)", m[1:len(m)-1])
	})

	re = regexp.MustCompile(`:[^/#?()\.\\]+`)
	regexpString = re.ReplaceAllStringFunc(regexpString, func(m string) string {
		return fmt.Sprintf(`(?P<%s>[^/#?]+)`, m[1:len(m)])
	})

	s := fmt.Sprintf(`\A%s\z`, regexpString)

	return &RegExRoute{
		AutoAck: false,
		Path:    path,
		Method:  method,
		Handler: fn,
		RegEx:   regexp.MustCompile(s),
	}
}

// Route represents a CoAP Route/Resource
type RegExRoute struct {
	Path       string
	Method     string
	Handler    RouteHandler
	RegEx      *regexp.Regexp
	AutoAck    bool
	MediaTypes []MediaType
}

func (r *RegExRoute) Matches(path string) (bool, map[string]string) {
	re := r.RegEx
	matches := re.FindAllStringSubmatch(path, -1)
	attrs := make(map[string]string)
	if len(matches) > 0 {
		subExp := re.SubexpNames()
		for idx, exp := range subExp {
			attrs[exp] = matches[0][idx]
		}
		return true, attrs
	}
	return false, attrs
}

func (r *RegExRoute) GetMethod() string {
	return r.Method
}

func (r *RegExRoute) GetMediaTypes() []MediaType {
	return r.MediaTypes
}

func (r *RegExRoute) GetConfiguredPath() string {
	return r.Path
}

func (r *RegExRoute) AutoAcknowledge() bool {
	return r.AutoAck
}

func (r *RegExRoute) Handle(req Request) Response {
	return r.Handler(req)
}

// MatchingRoute checks if a given path matches any defined routes/resources
func MatchingRoute(path string, method string, cf interface{}, routes []Route) (Route, map[string]string, error) {
	for _, route := range routes {
		if method == route.GetMethod() {
			match, attrs := route.Matches(path)

			if match {
				if len(route.GetMediaTypes()) > 0 {
					if cf == nil {
						return route, attrs, ErrUnsupportedContentFormat
					}

					foundMediaType := false
					for _, o := range route.GetMediaTypes() {
						if uint32(o) == cf {
							foundMediaType = true
							break
						}
					}

					if !foundMediaType {
						return route, attrs, ErrUnsupportedContentFormat
					}
				}
				return route, attrs, nil
			}
		}
	}
	return nil, nil, ErrNoMatchingRoute
}
