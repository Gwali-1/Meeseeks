package meeseeks

import (
	"net/http"
	"strings"
)

//struct http.Handler

type serverMux struct {
	NotFound         http.Handler
	MethodNotAllowed http.Handler
	registeredRoutes *[]route
	middlewares      []func(http.HandlerFunc) http.HandlerFunc
}

// struct to hold route handler info
type route struct {
	method  string
	path    []string
	handler http.Handler
}

//methods on it

// method to register handlers on patterns

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




func (s *serverMux) wrap(handler http.HandlerFunc)http.HandlerFunc{
	for _,m := range s.middlewares{
		handler = m(handler)
	}
	return handler
}




//method to wrap middlewares around handlers

//method serverHttp

//match function

//extract path parameter value

//custom
//middleware to server files from a directory
//middleware for jwt handling
//middleware protected routes by specifying middleware name
