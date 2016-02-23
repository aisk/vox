package gibson

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

var NotFound = &Response{
	StatusCode: 404,
	Content:    []byte("not found"),
}
