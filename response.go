package vox

import (
	"mime"
	"net/http"
	"regexp"
)

var bodyMatcher = regexp.MustCompile("^\\s*<")

// A Response is for all infomation, which will write to ResponseWriter.
type Response struct {
	body           interface{}
	explicitStatus bool
	status         int
	Header         http.Header
}

// Body is the response's body getter
func (response *Response) Body() interface{} {
	if response.body == nil {
		return "Not found"
	}
	return response.body
}

// SetBody is the response's body getter
func (response *Response) SetBody(body interface{}) {
	if body == nil {
		response.SetStatus(204)
	}
	if !response.explicitStatus {
		response.SetStatus(200)
	}
	switch v := body.(type) {
	case []byte:
		if bodyMatcher.Match(v) {
			response.Header.Set("Content-Type", mime.TypeByExtension(".html"))
		} else {
			response.Header.Set("Content-Type", mime.TypeByExtension(".text"))
		}
	case string:
		if bodyMatcher.MatchString(v) {
			response.Header.Set("Content-Type", mime.TypeByExtension(".html"))
		} else {
			response.Header.Set("Content-Type", mime.TypeByExtension(".text"))
		}
	case map[string]interface{}, map[string]string:
		response.Header.Set("Content-Type", mime.TypeByExtension(".json"))
	default:
		// TODO: suport more body types
	}

	response.body = body
}

// Status returns the response's status code
func (response *Response) Status() int {
	return response.status
}

// SetStatus is the response's status setter
func (response *Response) SetStatus(code int) {
	response.explicitStatus = true
	response.status = code
}

func createResponse(rw http.ResponseWriter) *Response {
	return &Response{
		status:         404,
		explicitStatus: false,
		Header:         rw.Header(),
	}
}
