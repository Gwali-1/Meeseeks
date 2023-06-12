package meeseeks

import (
	"net/http"
	"strings"
	"context"
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

func NewMeeseeks() *serverMux {
	return &serverMux{
		NotFound: http.NotFoundHandler(),
		MethodNotAllowed: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}),
		registeredRoutes: &[]route{},
	}
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

// method to wrap middlewares around handlers
func (s *serverMux) wrap(handler http.HandlerFunc) http.HandlerFunc {
	for _, m := range s.middlewares {
		handler = m(handler)
	}
	return handler
}

//method serverHttp

//match function
func (r route) match(url string, c context.ContextC){

}

//extract path parameter value

//custom
//middleware to server files from a directory
//middleware for jwt handling
//middleware protected routes by specifying middleware name
