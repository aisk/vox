package vox

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

// ErrNotAcceptable is the error returns when vox found the reqeust is not acceptable.
var ErrNotAcceptable = errors.New("content is not acceptable")

// A BaseRequest object contains all the information from current HTTP client.
//
// BaseRequest embedded the current request's raw *http.Request as it's field, so you
// can using all the fields and method of http.Request. see http://golang.org/pkg/net/http/#Request.
type BaseRequest struct {
	*http.Request

	// Params the parameters which extracted from the route.
	//
	// If the registered route is "/hello/{name}", and the actual path which
	// visited is "/hello/jim", the Params should be map[string]{"name": "jim"}
	//
	// Multiple parameters with same key is invalid and will be ignored.
	Params map[string]string

	response *BaseResponse
}

func createRequest(raw *http.Request) *BaseRequest {
	return &BaseRequest{
		raw,
		make(map[string]string),
		nil,
	}
}

// JSON is a helper to decode JSON request body to go value, with additional functionality to check the content type header from the request. If the content type header do not starts with "application/json" or decode errors, this function will return an error and set the response status code to 406.
func (request *BaseRequest) JSON(v interface{}) error {
	if !strings.HasPrefix(request.Header.Get("content-type"), "application/json") {
		request.response.Status = 406
		return ErrNotAcceptable
	}
	err := json.NewDecoder(request.Body).Decode(v)
	if err != nil {
		request.response.Status = 406
	}
	return err
}

// NoBody disables automatic request decoding. As an output type it produces an
// empty response, with status 204 unless the handler sets a status explicitly.
type NoBody struct{}

// Request contains a decoded body and shared HTTP request metadata.
// Use req.Request.Body to access the original stream.
type Request[T any] struct {
	*BaseRequest
	Body T
}
