package vox

import (
	"net/http"
)

// Context is a context through a HTTP request.
type Context struct {
	Req *Request
}

func createContext(rq *http.Request, rw http.ResponseWriter) *Context {
	return &Context{
		Req: createRequest(rq),
	}
}
