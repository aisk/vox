package vox

import (
	"regexp"
)

// Route will register a new path handler to a given path.
func (app *Application) Route(method string, path *regexp.Regexp, fn interface{}) {
	// TODO: support string based path
	switch v := fn.(type) {
	case func(req *Request, res *Response):
		app.middlewares = append(app.middlewares, func(req *Request, res *Response, next func()) {
			if match(req, method, path) {
				v(req, res)
				return
			}
			next()
		})
	case func(req *Request, res *Response, next func()):
		app.middlewares = append(app.middlewares, func(req *Request, res *Response, next func()) {
			if match(req, method, path) {
				v(req, res, next)
				return
			}
			next()
		})
	default:
		panic("invalid middleware function signature")
	}
}

func match(req *Request, method string, path *regexp.Regexp) bool {
	if req.Method != method {
		// TODO(asaka): ignore case?
		return false
	}
	match := path.FindStringSubmatch(req.URL.Path)
	if match == nil {
		return false
	}
	for i, name := range path.SubexpNames() {
		if i == 0 || name == "" {
			continue
		}
		req.Params[name] = match[i]
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
