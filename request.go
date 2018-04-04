package vox

import (
	"net/http"
	"net/url"
)

// A Request stores current HTTP request's infomation.
type Request struct {
	*http.Request
	Query  url.Values
	Params map[string]string
	Origin string
}

func createRequest(raw *http.Request) *Request {
	return &Request{
		raw,
		raw.URL.Query(),
		map[string]string{},
		raw.Header.Get("Origin"),
	}
}
