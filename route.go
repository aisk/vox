package vox

import (
	"regexp"
)

// Route will register a new path handler to a given path.
func (app *Application) Route(method string, path string, fn interface{}) {
	reg, err := pathToRegexp(path)
	if err != nil {
		panic("compile path error:" + err.Error())
	}

	switch v := fn.(type) {
	case func(*Context):
		app.middlewares = append(app.middlewares, func(ctx *Context, next func()) {
			if ctx.Request.Method == method && reg.MatchString(ctx.Request.URL.Path) {
				v(ctx)
				return
			}
			next()
		})
	case func(*Context, func()):
		app.middlewares = append(app.middlewares, func(ctx *Context, next func()) {
			if ctx.Request.Method == method && reg.MatchString(ctx.Request.URL.Path) {
				v(ctx, next)
				return
			}
			next()
		})
	default:
		panic("invalid middleware function signature")
	}
}

func pathToRegexp(path string) (*regexp.Regexp, error) {
	// TODO(asaka): support parameters in path
	return regexp.Compile(path)
}

// Get register a new path handler for GET method
func (app *Application) Get(path string, fn interface{}) {
	app.Route("GET", path, fn)
}

// Post register a new path handler for GET method
func (app *Application) Post(path string, fn interface{}) {
	app.Route("POST", path, fn)
}

// Put register a new path handler for GET method
func (app *Application) Put(path string, fn interface{}) {
	app.Route("PUT", path, fn)
}

// Delete register a new path handler for GET method
func (app *Application) Delete(path string, fn interface{}) {
	app.Route("DELETE", path, fn)
}

// Option register a new path handler for GET method
func (app *Application) Option(path string, fn interface{}) {
	app.Route("OPTION", path, fn)
}
