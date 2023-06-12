package meeseeks

import (
	"net/http"
	"strings"
	"context"
)

// struct http.Handler
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

// initialize a new meeseeks router
func NewMeeseeks() *serverMux {
	return &serverMux{
		NotFound: http.NotFoundHandler(),
		MethodNotAllowed: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}),
		registeredRoutes: &[]route{},
	}
}

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

// method to wrap middlewares around handlers
func (s *serverMux) wrap(handler http.HandlerFunc) http.HandlerFunc {
	for _, m := range s.middlewares {
		handler = m(handler)
	}
	return handler
}

// middleware function must have http.HandlerFunc function signature thus func(http.ResponseWriter, *http.Request)
func (s *serverMux) Use(m ...func(http.HandlerFunc) http.HandlerFunc) {
	s.middlewares = append(s.middlewares, m...)

}

// match function
func (r route) match(c context.Context, requestURL string) (context.Context, bool) {
	urlPathSegements := strings.Split(requestURL, "/")
	if len(r.path) != len(urlPathSegements) {
		return c, false
	}

	for i, segement := range r.path {
		if strings.HasPrefix(segement, ":") {
			if urlPathSegements[i] != "" {
				key := strings.Split(segement, ":")[1]
				c = context.WithValue(c, string(key), urlPathSegements[i])
				continue
			}
			return c, false
		}

		if segement != urlPathSegements[i] {
			return c, false
		}
	}

	return c, true
}

// extract path parameter value
func loadParam(c context.Context, paramName string) string {
	value, ok := c.Value(string(paramName)).(string) //type assertion
	//if value is not of type string
	if !ok {
		return ""
	}
	return value
}

//custom
//middleware to server files from a directory
//middleware for jwt handling
//middleware protected routes by specifying middleware name
// /home/:name/love   /home//love    /home/marco/love
