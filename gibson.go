package gibson

import (
	"net/http"
)

type Gibson struct {
	Handlers []Handler
}

func New() *Gibson {
	gibson := &Gibson{
		Handlers: []Handler{
			NotFoundHandler,
		},
	}

	return gibson
}

func (gibson *Gibson) ServeHTTP(rw http.ResponseWriter, rq *http.Request) {
	var res *Response
	req := NewRequest(rw, rq)
	for i := len(gibson.Handlers) - 1; i >= 0; i-- {
		handler := gibson.Handlers[i]
		res = handler(req)
		if res != nil {
			break
		}
	}
	res.write(rw)
}

func (gibson *Gibson) Use(handler Handler) {
	gibson.Handlers = append(gibson.Handlers, handler)
}

func (gibson *Gibson) Route() {

}

func (gibson *Gibson) Run(addr string) {
	http.ListenAndServe(addr, gibson)
}
