<img src="https://upload.wikimedia.org/wikipedia/en/1/1d/Mr._Meeseeks.png" alt="drawing" width="100"/>

# Meeseeks

Meeseeks is lightweigth and provides powerful http routing functionality

#### current meeseeks features

- you can have **named parameters** in your routes.
- you can provide **Custom handlers** for `404 Not Found` and `405 Method Not Allowed` responses.
- Sets an `Allow` header for all `405 Method Not Allowed` responses.
- create and register middlewares to be used on all your routes
- Provide `http.HandlerFunc` functions as your path handlers and middlewares(strightfoward and easy)
- Zero dependencies.
- lightweight and gives you maximum control

### Installation

```
$ go get github.com/Gwali-1/Meeseeks
```

### Example Usage

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    " github.com/Gwali-1/Meeseeks"

func main() {
    // Initialize a new meeseeks router.
    router := meeseeks.NewMeeseeks()

    //register a handler function to handle get request on the /rocket/:name route
    //you can use : to denote a named path parameter
    //register handlers with route patterns using the http verb named methods on the meeseeks router (currently supports just GET and POST)
    router.GET("/rocket/:name", func(w http.ResponseWriter, r *http.Request) {
    // Using  meeseeks.LoadParam()  you can retrieve the path parameter specifed in you route with by passing
    //the request context and the parameter name
    // request context.
    name := meeseeks.LoadParam(r.Context(), "name")
    fmt.Fprintf(w, "rocket name is  %s", name)
})

    err := http.ListenAndServe(":2323", mux)
    log.Fatal(err)
}
```

### Using Middleware

```go
 router := meeseeks.NewMeeseeks()

// The Use() method can be used to register middleware. Middleware registered
//on the router will be available on routes registered after the after the middle are regitration.
mux.Use(exampleMiddleware1)
mux.Use(exampleMiddleware2)
mux.Use(exampleMiddleware3)
mux.Use(exampleMiddleware4)

//middleware excution is bottom to top meaning exampleMiddleware4 is executed first in the middleware chain then
//exampleMiddleware3 and so on. This is useful if you need your middleware to execute in a particular order

```

### Things To Note

- You can have conflicting routes like (e.g. `/posts/:id` and `posts/new`). The order in which they are declared is how they are matched
- Trailing slashes are significant ; `/profile/:id` and `/profile/:id/` are not the same and won't match.
- An `Allow` header is automatically set for all `405 Method Not Allowed` responses
- Middleware must be declared _before_ a route in order to be used by that route. Any middleware declared after a route won't act on that route.
- For example:

```go
router := meeseeks.NewMeeseeks()
router.Use(middleware1)
router.GET("/foo",foo) // This route will use middleware1 only.
router.Use(middleware2)
router.GET("/bar",bar)  // This route will use both middleware1 and middleware2.
```

### Contributing

Feel free to send in any contibutions (bug fixes , feature sugestions , etc)

Was motivated to write meeseeks after reading source code for  [alexedwards/flow](https://github.com/alexedwards/flow)

### Todo
- in memory session store
- provide in built middlewares for file serving
- middleware tests

