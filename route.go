package vox

import (
	"regexp"
)

// Route will register a new path handler to a given path.
func (app *Application) Route(method string, path *regexp.Regexp, handler Handler) {
	// TODO: support string based path
	app.middlewares = append(app.middlewares, func(req *Request, res *Response) {
		if match(req, method, path) {
			handler(req, res)
			return
		}
		req.Next()
	})
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
func (app *Application) Get(path *regexp.Regexp, handler Handler) {
	app.Route("GET", path, handler)
}

// Post register a new path handler for GET method
func (app *Application) Post(path *regexp.Regexp, handler Handler) {
	app.Route("POST", path, handler)
}

// Put register a new path handler for GET method
func (app *Application) Put(path *regexp.Regexp, handler Handler) {
	app.Route("PUT", path, handler)
}

// Delete register a new path handler for GET method
func (app *Application) Delete(path *regexp.Regexp, handler Handler) {
	app.Route("DELETE", path, handler)
}

// Option register a new path handler for GET method
func (app *Application) Option(path *regexp.Regexp, handler Handler) {
	app.Route("OPTION", path, handler)
}
