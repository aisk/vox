package vox

import (
	"regexp"
)

// Route will register a new path handler to a given path.
func (app *Application) Route(method string, path string, handler Handler) {
	// TODO: support string based path
	r := routeToRegexp(path)
	app.middlewares = append(app.middlewares, func(req *Request, res *Response) {
		if match(req, method, r) {
			handler(req, res)
			return
		}
		req.Next()
	})
}

func match(req *Request, method string, path *regexp.Regexp) bool {
	if req.Method != method && method != "*" {
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

// Get register a new path handler for GET method.
func (app *Application) Get(path string, handler Handler) {
	app.Route("GET", path, handler)
}

// Head register a new path handler for HEAD method.
func (app *Application) Head(path string, handler Handler) {
	app.Route("HEAD", path, handler)
}

// Post register a new path handler for POST method.
func (app *Application) Post(path string, handler Handler) {
	app.Route("POST", path, handler)
}

// Put register a new path handler for PUT method.
func (app *Application) Put(path string, handler Handler) {
	app.Route("PUT", path, handler)
}

// Patch register a new path handler for PATCH method.
func (app *Application) Patch(path string, handler Handler) {
	app.Route("PATCH", path, handler)
}

// Delete register a new path handler for DELETE method.
func (app *Application) Delete(path string, handler Handler) {
	app.Route("DELETE", path, handler)
}

// Connect register a new path handler for CONNECT method.
func (app *Application) Connect(path string, handler Handler) {
	app.Route("CONNECT", path, handler)
}

// Options register a new path handler for OPTIONS method.
func (app *Application) Options(path string, handler Handler) {
	app.Route("OPTIONS", path, handler)
}

// Trace register a new path handler for TRACE method.
func (app *Application) Trace(path string, handler Handler) {
	app.Route("TRACE", path, handler)
}

func routeToRegexp(path string) *regexp.Regexp {
	replaced := regexp.MustCompile(`{(?P<param>\w+)}`).ReplaceAllStringFunc(path, func(s string) string {
		return "(?P<" + s[1:len(s)-1] + ">\\w+)"
	})
	return regexp.MustCompile("^" + replaced + "$")
}
