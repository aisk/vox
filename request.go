package vox

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strings"
)

// A DecodeError describes why a JSON request body was rejected. Its message is
// safe to send to the client; Err holds the underlying error, if any.
type DecodeError struct {
	// Status is the HTTP status code the failure maps to.
	Status  int
	Message string
	Err     error
}

func (e *DecodeError) Error() string { return e.Message }

func (e *DecodeError) Unwrap() error { return e.Err }

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

	app      *Application
	response *BaseResponse
}

func createRequest(raw *http.Request) *BaseRequest {
	return &BaseRequest{
		Request: raw,
		Params:  make(map[string]string),
	}
}

// JSON decodes the JSON request body into v. Route handlers with an input type
// other than NoBody have it called for them. On failure it sets the response
// status and body to a *DecodeError and returns it.
func (request *BaseRequest) JSON(v any) error {
	if err := request.decodeJSON(v); err != nil {
		request.response.Status = err.Status
		request.response.Body = err
		return err
	}
	return nil
}

// decodeJSON reads exactly one JSON value other than null from the request
// body into v.
func (request *BaseRequest) decodeJSON(v any) *DecodeError {
	mediaType, _, err := mime.ParseMediaType(request.Header.Get("Content-Type"))
	if err != nil || (mediaType != "application/json" &&
		!(strings.HasPrefix(mediaType, "application/") && strings.HasSuffix(mediaType, "+json"))) {
		return &DecodeError{http.StatusUnsupportedMediaType, "content type must be application/json", err}
	}

	body := request.Body
	if limit := request.app.maxBodySize; limit > 0 {
		body = http.MaxBytesReader(request.response.Writer, body, limit)
	}
	data, err := io.ReadAll(body)
	if err != nil {
		var tooLarge *http.MaxBytesError
		if errors.As(err, &tooLarge) {
			return &DecodeError{http.StatusRequestEntityTooLarge, "request body is too large", err}
		}
		return &DecodeError{http.StatusBadRequest, "request body could not be read", err}
	}

	switch trimmed := bytes.TrimSpace(data); {
	case len(trimmed) == 0:
		return &DecodeError{http.StatusBadRequest, "request body is empty", nil}
	case string(trimmed) == "null":
		return &DecodeError{http.StatusBadRequest, "request body must not be null", nil}
	}

	if err := json.Unmarshal(data, v); err != nil {
		message := "request body is not valid JSON"
		var syntaxErr *json.SyntaxError
		var typeErr *json.UnmarshalTypeError
		switch {
		case errors.As(err, &syntaxErr):
			message = fmt.Sprintf("malformed JSON at offset %d", syntaxErr.Offset)
		case errors.As(err, &typeErr) && typeErr.Field != "":
			message = fmt.Sprintf("invalid type for field %q", typeErr.Field)
		case errors.As(err, &typeErr):
			message = "invalid type for request body"
		}
		return &DecodeError{http.StatusBadRequest, message, err}
	}
	return nil
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
