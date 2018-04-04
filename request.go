package vox

import (
	"net/http"
)

// A Request stores current HTTP request's infomation.
type Request struct {
	*http.Request
	Params map[string]string
	State  map[string]interface{}
}

func createRequest(raw *http.Request) *Request {
	return &Request{
		raw,
		make(map[string]string),
		make(map[string]interface{}),
	}
}
