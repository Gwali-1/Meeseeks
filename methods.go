package meeseeks

import (
	"net/http"
	"strings"
)


//this file contains all methods on the serveMux struct used to register patterns
//with handlers
//they are specified with http verbs indicating the request method to be used in the request
//fucntions with http.HandlerFunc type are accepted as the argument to be used as request handler for conveninece. They are converted to handlers later on
//the whole router interface is designed to work with functions of type http.HandlerFunc from the users end

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
