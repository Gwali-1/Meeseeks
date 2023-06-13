package meeseeks

import (
	"context"
	"net/http"
	"strings"
)

//the serverMux type is a http.Handler that contains other handlers
//it implementation of the ServeHTTP function on the http.Handler interface makes it so
//the type has various metohds implemented for registering handlers with url paths and also wrapping handlers with middleware to be executed within the request response cylce
//there is another struct of type router that contains url patterns , thier respective handlers etc
//it contains methods to match incoming request url against a other registered patterns and find a match
//the match function contain the logic for matching the requested url to its appropriate handler. the matching function is kept minimal here and is higly influenced by the matching
//function in github.com/alexedwards/flow

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

// serverMux implements serveHTTP method hence is of type http.Handler
func (s *serverMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	methodsAllowed := []string{}
	for _, route := range *s.registeredRoutes {
		ctx, match := route.match(r.Context(), r.URL.Path)
		if match {
			if route.method == r.Method {
				route.handler.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			if !contains(methodsAllowed, route.method) {
				methodsAllowed = append(methodsAllowed, route.method)
			}
		}
	}

	if len(methodsAllowed) > 0 {
		w.Header().Set("Allow", strings.Join(append(methodsAllowed, http.MethodOptions), ", "))
		f := http.HandlerFunc(s.MethodNotAllowed.ServeHTTP) //convert the http.Handler to http.HandleFunc so we can wrap middlewares
		s.wrap(f).ServeHTTP(w, r)
		return
	}
	f := http.HandlerFunc(s.NotFound.ServeHTTP) //convert the http.Handler to http.HandleFunc so we can wrap middlewares
	s.wrap(f).ServeHTTP(w, r)

}

// request url  match function
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

// extract path parameter value by passing in request context
func LoadParam(c context.Context, paramName string) string {
	value, ok := c.Value(string(paramName)).(string) //type assertion
	//if value is not of type string
	if !ok {
		return ""
	}
	return value
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
