package vox

type Handler func(*Request) *Response

func NotFoundHandler(req *Request) *Response {
	return NotFound()
}
