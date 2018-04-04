package vox

import (
	"encoding/json"
	"net/http"
)

// Application is type of an vox application.
type Application struct {
	middlewares []func(*Request, *Response, func())
	fn          func(*Request, *Response)
}

// New returns a new vox Application.
func New() *Application {
	app := &Application{}
	return app
}

// Use a vox middleware.
func (app *Application) Use(fn interface{}) {
	switch v := fn.(type) {
	case func(*Request, *Response):
		app.middlewares = append(app.middlewares, func(req *Request, res *Response, _ func()) {
			v(req, res)
		})
	case func(*Request, *Response, func()):
		app.middlewares = append(app.middlewares, func(req *Request, res *Response, next func()) {
			v(req, res, next)
		})
	default:
		panic("invalid middleware function signature")
	}
}

func (app *Application) ServeHTTP(rw http.ResponseWriter, rq *http.Request) {
	app.fn = compose(app.middlewares)
	req := createRequest(rq)
	res := createResponse(rw)
	app.fn(req, res)
	respond(res)
}

// Run the Vox application.
func (app *Application) Run(addr string) {
	http.ListenAndServe(addr, app)
}

func compose(middlewares []func(*Request, *Response, func())) func(*Request, *Response) {
	return func(req *Request, res *Response) {
		length := len(middlewares)
		nexts := make([]func(), length+1)
		nexts[length] = func() {}
		for i := length; i > 0; i-- {
			func(j int) {
				nexts[j-1] = func() {
					middlewares[j-1](req, res, nexts[j])
				}
			}(i)
		}
		nexts[0]()
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
