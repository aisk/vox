package vox

import (
	"net/http"
)

type Vox struct {
	MiddleWares []MiddleWare
	Handlers    []Handler
}

func New() *Vox {
	vox := &Vox{
		MiddleWares: []MiddleWare{},
		Handlers: []Handler{
			NotFoundHandler,
		},
	}
	return vox
}

func (vox *Vox) Use(middleware MiddleWare) {
	vox.MiddleWares = append(vox.MiddleWares, middleware)
}

func (vox *Vox) ServeHTTP(rw http.ResponseWriter, rq *http.Request) {
	req := newRequest(rw, rq)
	res := NotFound()
	middlewareIndex := 0

	var next func() *Response
	next = func() *Response {
		if middlewareIndex == len(vox.MiddleWares) {
			return res
		}
		middlewareIndex++
		res = vox.MiddleWares[middlewareIndex-1](req, next)
		return res
	}

	next()
	res.write(rw)
}

func (vox *Vox) Route(path string, handler Handler) {
}

func (vox *Vox) Run(addr string) {
	http.ListenAndServe(addr, vox)
}
