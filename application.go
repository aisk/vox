package vox

import (
	"encoding/json"
	"net/http"
)

// Application is type of an vox application.
type Application struct {
	middlewares []Handler
}

// New returns a new vox Application.
func New() *Application {
	app := &Application{}
	return app
}

// Use a vox middleware.
func (app *Application) Use(handler Handler) {
	app.middlewares = append(app.middlewares, handler)
}

func (app *Application) ServeHTTP(rw http.ResponseWriter, rq *http.Request) {
	handler := compose(app.middlewares)
	req := createRequest(rq)
	res := createResponse(rw)
	handler(req, res)
	respond(res)
}

// Run the Vox application.
func (app *Application) Run(addr string) error {
	return http.ListenAndServe(addr, app)
}

func compose(middlewares []Handler) Handler {
	return func(req *Request, res *Response) {
		next := func() {}
		for i := len(middlewares) - 1; i >= 0; i-- {
			func(i int, nenext func()) {
				next = func() {
					req.Next = nenext
					middlewares[i](req, res)
				}
			}(i, next)
		}
		next()
	}
}

func respond(res *Response) {
	res.setImplict()

	res.Writer.WriteHeader(res.Status)

	switch v := res.Body.(type) {
	case []byte:
		res.Writer.Write(v)
	case string:
		res.Writer.Write([]byte(v))
	// case map[string]string, map[string]interface{}:
	default:
		body, err := json.Marshal(v)
		if err != nil {
			panic(err)
		}
		res.Writer.Write(body)
	}
}
