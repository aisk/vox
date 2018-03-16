package vox

import (
	"context"
	"net/http"
)

// Context is a context through a HTTP request.
type Context struct {
	Context  context.Context
	Request  *Request
	Response *Response
	App      *Application           // Vox application instance
	Req      *http.Request          // Origin go's http.Request for this request
	Res      http.ResponseWriter    // Origin go's http.ResponseWriter for this request
	State    map[string]interface{} // Recommended place to store custom data for middlewares
	Cookies  *Cookie
}

func (app *Application) createContext(rq *http.Request, rw http.ResponseWriter) *Context {
	ctx := &Context{
		Context:  rq.Context(),
		Request:  createRequest(rq),
		Response: createResponse(rw),
		App:      app,
		Req:      rq,
		Res:      rw,
		State:    map[string]interface{}{},
		Cookies:  createCookies(rq, rw),
	}
	return ctx
}
