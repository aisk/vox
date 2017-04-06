package vox

import (
	"encoding/json"
	"net/http"
)

// Application is type of an vox application.
type Application struct {
	middlewares []func(*Context, func())
	fn          func(*Context)
}

// New returns a new vox Application.
func New() *Application {
	app := &Application{}
	return app
}

// Use a vox middleware.
func (app *Application) Use(fn interface{}) {
	switch v := fn.(type) {
	case func(*Context):
		app.middlewares = append(app.middlewares, func(ctx *Context, _ func()) {
			v(ctx)
		})
	case func(*Context, func()):
		app.middlewares = append(app.middlewares, func(ctx *Context, next func()) {
			v(ctx, next)
		})
	default:
		panic("invalid middleware function signature")
	}
}

func (app *Application) ServeHTTP(rw http.ResponseWriter, rq *http.Request) {
	app.fn = compose(app.middlewares)
	ctx := app.createContext(rq, rw)
	app.fn(ctx)
	respond(ctx)
}

func (app *Application) createContext(rq *http.Request, rw http.ResponseWriter) *Context {
	ctx := &Context{
		Request:  createRequest(rq),
		Response: createResponse(rw),
		Req:      rq,
		Res:      rw,
		App:      app,
	}
	return ctx
}

// Route will register a new path handler to a given path.
func (app *Application) Route(path string) {
}

// Run the Vox application.
func (app *Application) Run(addr string) {
	http.ListenAndServe(addr, app)
}

func compose(middlewares []func(*Context, func())) func(*Context) {
	return func(ctx *Context) {
		length := len(middlewares)
		nexts := make([]func(), length+1)
		nexts[length] = func() {}
		for i := length; i > 0; i-- {
			func(j int) {
				nexts[j-1] = func() {
					middlewares[j-1](ctx, nexts[j])
				}
			}(i)
		}
		nexts[0]()
	}
}

func respond(ctx *Context) {
	ctx.Res.WriteHeader(ctx.Response.Status)
	switch v := ctx.Response.Body.(type) {
	case []byte:
		ctx.Res.Write(v)
	case string:
		ctx.Res.Write([]byte(v))
	case nil:
		ctx.Res.Write(nil)
	case map[interface{}]interface{}:
		body, err := json.Marshal(v)
		if err != nil {
			panic(err)
		}
		ctx.Res.Write(body)
	default:
		panic("not implemented yed")
	}
}
