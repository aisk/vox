package vox

import (
	"net/http"
)

type Response struct {
	StatusCode int
	Content    []byte
	Header     map[string]string
}

func (res *Response) write(rw http.ResponseWriter) error {
	for name, value := range res.Header {
		rw.Header().Set(name, value)
	}
	rw.WriteHeader(res.StatusCode)
	_, err := rw.Write(res.Content)
	return err
}

func NewResponse(content interface{}, args ...interface{}) *Response {
	res := &Response{StatusCode: 200}

	switch c := content.(type) {
	case string:
		res.Content = []byte(c)
	case []byte:
		res.Content = c
	case interface {
		String() string
	}:
		res.Content = []byte(c.String())
	}

	if len(args) == 1 {
		switch arg := args[0].(type) {
		case int:
			res.StatusCode = arg
		case map[string]string:
			res.Header = arg
		}
	}

	if len(args) == 2 {
		if statusCode, ok := args[0].(int); ok {
			res.StatusCode = statusCode
		} else if statusCode, ok := args[1].(int); ok {
			res.StatusCode = statusCode
		}
		if header, ok := args[0].(map[string]string); ok {
			res.Header = header
		} else if header, ok := args[1].(map[string]string); ok {
			res.Header = header
		}

	}

	return res
}

func NotFound(args ...string) *Response {
	res := &Response{
		StatusCode: 404,
	}
	if len(args) == 1 {
		res.Content = []byte("not found")
	}
	return res
}

func JSON() {

}

func Text(content string, args ...int) *Response {
	res := &Response{
		StatusCode: 200,
		Content:    []byte(content),
	}
	if len(args) == 1 {
		res.StatusCode = args[0]
	}
	return res
}
