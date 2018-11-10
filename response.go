package vox

import (
	"mime"
	"net/http"
	"regexp"
)

var (
	explictSetStatus = -1
	explictSetBody   = struct{}{}

	bodyMatcher = regexp.MustCompile("^\\s*<")
)

// A Response object contains all the information which will written to current
// HTTP client.
type Response struct {
	// Writer is the raw http.ResponseWriter for current request. You should
	// assign the Body / Status / Header value instead of using this field.
	Writer http.ResponseWriter
	// Body is the container for HTTP response's body.
	Body interface{}
	// The status code which will respond as the HTTP Response's status code.
	// 200 will be used as the default value if not set.
	Status int
	// Headers which will be written to the response.
	Header http.Header
}

func (response *Response) setImplictStatus() {
	if response.Status != explictSetStatus {
		// response's status is set by user.
		return
	}
	if response.Body == explictSetBody {
		// response's body is not set, set it to 404.
		response.Status = 404
		return
	}
	// response's status is not set by user, give it a default value.
	response.Status = 200
}

func (response *Response) setImplictContentType() {
	if response.Header.Get("Content-Type") != "" {
		return
	}

	if response.Body == explictSetBody {
		return
	}

	switch v := response.Body.(type) {
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
	// case map[string]interface{}, map[string]string:
	default:
		response.Header.Set("Content-Type", mime.TypeByExtension(".json"))
	}
}

func (response *Response) setImplictBody() {
	if response.Body == explictSetBody {
		response.Body = ""
	}
}

func (response *Response) setImplict() {
	response.setImplictStatus()
	response.setImplictContentType()
	response.setImplictBody()
}

func createResponse(rw http.ResponseWriter) *Response {
	return &Response{
		Writer: rw,
		Body:   explictSetBody,
		Status: explictSetStatus,
		Header: rw.Header(),
	}
}
