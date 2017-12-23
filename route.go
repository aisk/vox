package vox

import (
	"regexp"
)

// Route will register a new path handler to a given path.
func (app *Application) Route(method string, path *regexp.Regexp, fn interface{}) {
	// TODO: support string based path
	switch v := fn.(type) {
	case func(*Context):
		app.middlewares = append(app.middlewares, func(ctx *Context, next func()) {
			if match(ctx, method, path) {
				v(ctx)
				return
			}
			next()
		})
	case func(*Context, func()):
		app.middlewares = append(app.middlewares, func(ctx *Context, next func()) {
			if match(ctx, method, path) {
				v(ctx, next)
				return
			}
			next()
		})
	default:
		panic("invalid middleware function signature")
	}
}

func match(ctx *Context, method string, path *regexp.Regexp) bool {
	if ctx.Request.Method != method {
		// TODO(asaka): ignore case?
		return false
	}
	match := path.FindStringSubmatch(ctx.Request.URL.Path)
	if match == nil {
		return false
	}
	for i, name := range path.SubexpNames() {
		if i == 0 || name == "" {
			continue
		}
		ctx.Request.Params[name] = match[i]
	}
	return true
}

// Get register a new path handler for GET method
func (app *Application) Get(path *regexp.Regexp, fn interface{}) {
	app.Route("GET", path, fn)
}

// Post register a new path handler for GET method
func (app *Application) Post(path *regexp.Regexp, fn interface{}) {
	app.Route("POST", path, fn)
}

// Put register a new path handler for GET method
func (app *Application) Put(path *regexp.Regexp, fn interface{}) {
	app.Route("PUT", path, fn)
}

// Delete register a new path handler for GET method
func (app *Application) Delete(path *regexp.Regexp, fn interface{}) {
	app.Route("DELETE", path, fn)
}

// Option register a new path handler for GET method
func (app *Application) Option(path *regexp.Regexp, fn interface{}) {
	app.Route("OPTION", path, fn)
}
