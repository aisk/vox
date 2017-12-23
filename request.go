package vox

import (
	"net/http"
	"net/url"
)

// A Request stores current HTTP request's infomation.
type Request struct {
	Request *http.Request
	Header  http.Header
	Method  string
	URL     *url.URL
	Params  map[string]string
}

func createRequest(raw *http.Request) *Request {
	return &Request{
		Request: raw,
		Header:  raw.Header,
		Method:  raw.Method,
		URL:     raw.URL,
		Params:  map[string]string{},
	}
}
