package vox

import (
	"net/http"
)

// Context is a context through a HTTP request.
type Context struct {
	Request  *Request
	Response *Response
	App      *Application
	Req      *http.Request
	Res      http.ResponseWriter
}
