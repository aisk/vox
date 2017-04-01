package vox

import (
	"net/http"
)

// A Response is for all infomation, which will write to ResponseWriter.
type Response struct {
	Status  int
	Body    interface{}
	Headers http.Header
}
