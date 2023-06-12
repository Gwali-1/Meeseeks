package meeseeks

import "net/http"

//struct http.Handler

type serverMux struct {
	NotFound         http.Handler
	MethodNotAllowed http.Handler
	registeredRoutes []*routes
	middlewares      []func(http.HandlerFunc) http.HandlerFunc
}

// struct to hold route handler info
type routes struct {
	method  string
	path    string
	handler http.Handler
}

//methods on it

// method to register handlers on patterns

//method to wrap middlewares around handlers

//method serverHttp

//match function

//extract path parameter value

//custom
//middleware to server files from a directory
//middleware for jwt handling
//middleware protected routes by specifying middleware name
