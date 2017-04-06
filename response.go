package vox

import (
	"net/http"
)

// A Response is for all infomation, which will write to ResponseWriter.
type Response struct {
	Status int
	Body   interface{}
	Header http.Header
}

func createResponse(rw http.ResponseWriter) *Response {
	return &Response{
		Status: 404,
		Header: rw.Header(),
	}
}
