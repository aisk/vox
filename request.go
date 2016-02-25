package vox

import (
	"net/http"
	"net/url"
)

type Request struct {
	RawRequest        *http.Request
	RawResponseWriter http.ResponseWriter

	// stored the context info from middlewares
	Context map[string]interface{}

	Header http.Header
	Method string
	URL    *url.URL
}

func NewRequest(rw http.ResponseWriter, rq *http.Request) *Request {
	req := &Request{
		RawRequest:        rq,
		RawResponseWriter: rw,

		Method: rq.Method,
		URL:    rq.URL,
		Header: rq.Header,
	}
	return req
}
