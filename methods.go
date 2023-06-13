package meeseeks

import (
	"net/http"
	"strings"
)

// serverMux struct methods
func (s *serverMux) GET(pattern string, handler http.HandlerFunc) {
	allowedMethods := []string{http.MethodGet, http.MethodHead}
	for _, method := range allowedMethods {
		newRoute := route{
			method:  method,
			path:    strings.Split(pattern, "/"),
			handler: http.HandlerFunc(s.wrap(handler)),
		}

		*s.registeredRoutes = append(*s.registeredRoutes, newRoute)
	}
}

func (s *serverMux) POST(pattern string, handler http.HandlerFunc) {
	allowedMethods := []string{http.MethodPost}
	for _, method := range allowedMethods {
		newRoute := route{
			method:  method,
			path:    strings.Split(pattern, "/"),
			handler: http.HandlerFunc(handler),
		}
		*s.registeredRoutes = append(*s.registeredRoutes, newRoute)
	}
}
