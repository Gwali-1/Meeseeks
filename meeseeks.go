package meeseeks

import (
	"context"
	"net/http"
	"strings"
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

//serverMux implements serveHTTP method hence is a http.Handler

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

