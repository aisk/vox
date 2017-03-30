package vox

import (
	"net/http"
)

// Application is type of an vox application.
type Application struct {
	middlewares []func(*Context, func())
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
	ctx := createContext(rq, rw)
	for _, middleware := range app.middlewares {
		middleware(ctx, func() {})
	}
	rw.Write([]byte("ok"))
}

// Route will register a new path handler to a given path.
func (app *Application) Route(path string) {
}

// Run the Vox application.
func (app *Application) Run(addr string) {
	http.ListenAndServe(addr, app)
}
