package vox

import (
	"net/http"
)

// A Request object contains all the information from current HTTP client.
//
// Request embedded the current request's raw *http.Request as it's field, so you
// can using all the fields and method of http.Request. see http://golang.org/pkg/net/http/#Request.
type Request struct {
	*http.Request
	// Params the parameters which extracted from the route.
	//
	// If the registered route is "/hello/{name}", and the actual path which
	// visited is "/hello/jim", the Params should be map[string]{"name": "jim"}
	//
	// Multiple parameters with same key is invalid and will be ignored.
	Params map[string]string
	// State is a place which is recommended as storage for passing information
	// through middleware and to your frontend views.
	//
	// State will be empty if you or your middlewares don't set it manually.
	//
	// It's recommended to using a namespace prefix to avoid key name conflict,
	// like "mysession:session_id", especially for middleware libraries.
	//
	// State's value type is just interface{}, so you can cast it to it's own
	// type using type assesions (https://tour.golang.org/methods/15).
	//
	// For middleware library authors, you should provide utility functions to
	// extract the actually value, with type casting and namespace demangling.
	// For example, mysession middleware can provide a function with type
	// mysesion.GetSession(*vox.Request) *mysession.Session
	State map[string]interface{}
	// Next will call the next handler / middleware to processing request.
	// It's the middleware's responsibility to call the Next function (or not).
	Next func()
}

func createRequest(raw *http.Request) *Request {
	return &Request{
		raw,
		make(map[string]string),
		make(map[string]interface{}),
		nil,
	}
}
