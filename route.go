package vox

// routeHandler handles route matching and parameter extraction
func (app *Application) routeHandler(ctx *Context, req *Request, res *Response) {
	match, found := app.router.Match(req.Method, "", req.URL.Path)
	if found {
		for k, v := range match.Params {
			req.Params[k] = v
		}
		h := match.Handler
		h(ctx, req, res)
	}
	ctx.Next()
}

// Route will register a new path handler to a given path.
func (app *Application) Route(method string, path string, handler Handler) {
	var err error
	if method == "*" {
		err = app.router.Handle(path, handler)
	} else {
		err = app.router.Handle(method+" "+path, handler)
	}
	if err != nil {
		panic(err)
	}
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

// Options register a new path handler for OPTIONS method.
func (app *Application) Options(path string, handler Handler) {
	app.Route("OPTIONS", path, handler)
}

// Trace register a new path handler for TRACE method.
func (app *Application) Trace(path string, handler Handler) {
	app.Route("TRACE", path, handler)
}
