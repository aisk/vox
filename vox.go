package vox

import (
	"net/http"
)

type Vox struct {
	Handlers []Handler
}

func New() *Vox {
	vox := &Vox{
		Handlers: []Handler{
			NotFoundHandler,
		},
	}

	return vox
}

func (vox *Vox) ServeHTTP(rw http.ResponseWriter, rq *http.Request) {
	var res *Response
	req := NewRequest(rw, rq)
	for i := len(vox.Handlers) - 1; i >= 0; i-- {
		handler := vox.Handlers[i]
		res = handler(req)
		if res != nil {
			break
		}
	}
	res.write(rw)
}

func (vox *Vox) Use(handler Handler) {
	vox.Handlers = append(vox.Handlers, handler)
}

func (vox *Vox) Route() {

}

func (vox *Vox) Run(addr string) {
	http.ListenAndServe(addr, vox)
}
