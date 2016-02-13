package canopus

import (
	"fmt"
	"regexp"
)

// CreateCompilableRoutePath creates a RegEx for a valid route path
func CreateCompilableRoutePath(route string) (*regexp.Regexp, bool) {
	var re *regexp.Regexp
	var isStatic bool

	regexpString := route

	isStaticRegexp := regexp.MustCompile(`[\(\)\?\<\>:]`)
	if !isStaticRegexp.MatchString(route) {
		isStatic = true
	}

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

	return regexp.MustCompile(s), isStatic
}

// CreateNewRoute creates a new Route object
func CreateNewRoute(path string, method string, fn RouteHandler) *Route {
	re, _ := CreateCompilableRoutePath(path)

	return &Route{
		AutoAck: false,
		Path:    path,
		Method:  method,
		Handler: fn,
		RegEx:   re,
	}
}

// MatchesRoutePath checks if a given path matches a regex route
func MatchesRoutePath(path string, re *regexp.Regexp) (bool, map[string]string) {
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

// Route represents a CoAP Route/Resource
type Route struct {
	Path       string
	Method     string
	Handler    RouteHandler
	RegEx      *regexp.Regexp
	AutoAck    bool
	MediaTypes []MediaType
}

// MatchingRoute checks if a given path matches any defined routes/resources
func MatchingRoute(path string, method string, cf interface{}, routes []*Route) (*Route, map[string]string, error) {
	for _, route := range routes {
		if method == route.Method {
			match, attrs := MatchesRoutePath(path, route.RegEx)

			if match {
				if len(route.MediaTypes) > 0 {
					if cf == nil {
						return route, attrs, ErrUnsupportedContentFormat
					}

					foundMediaType := false
					for _, o := range route.MediaTypes {
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
