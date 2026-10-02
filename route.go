package vox

import (
	"encoding/json"
	"io"
	"mime"
	"net/http"
	"reflect"
	"strings"
)

// routeHandler handles route matching and parameter extraction
func (app *Application) routeHandler(ctx *Context, req *BaseRequest, res *BaseResponse) {
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

// registerRoute stores a uniform handler in the route tree.
func (app *Application) registerRoute(method string, path string, handler Handler) {
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

// Route registers a route handler. Inputs other than NoBody are decoded as JSON.
// Route handlers return to continue processing; Context.Next is for middleware only.
func (app *Application) Route[In, Out any](method, path string, handler RouteHandler[In, Out]) {
	app.registerRoute(method, path, func(ctx *Context, req *BaseRequest, res *BaseResponse) {
		var body In
		if _, skip := any(body).(NoBody); !skip {
			mediaType, _, err := mime.ParseMediaType(req.Header.Get("Content-Type"))
			if err != nil || (mediaType != "application/json" &&
				!(strings.HasPrefix(mediaType, "application/") && strings.HasSuffix(mediaType, "+json"))) {
				res.Status = http.StatusUnsupportedMediaType
				res.Body = http.StatusText(res.Status)
				return
			}
			decoder := json.NewDecoder(req.Body)
			if err := decoder.Decode(&body); err != nil {
				res.Status = http.StatusBadRequest
				res.Body = http.StatusText(res.Status)
				return
			}
			// Require exactly one JSON value, allowing trailing whitespace.
			var extra any
			if err := decoder.Decode(&extra); err != io.EOF {
				res.Status = http.StatusBadRequest
				res.Body = http.StatusText(res.Status)
				return
			}
		}
		request := &Request[In]{BaseRequest: req, Body: body}
		response := &Response[Out]{BaseResponse: res}
		next := ctx.Next
		defer func() { ctx.Next = next }()
		ctx.Next = func() {
			panic("vox: Context.Next is only available in middleware; return from a route handler instead")
		}
		handler(ctx, request, response)
		// A body written to BaseResponse directly (error, stream) takes precedence.
		if res.DontRespond || res.redirected || res.HasBody() {
			return
		}
		if _, empty := any(response.Body).(NoBody); empty {
			res.Body = []byte{}
			if res.Status == 0 {
				res.Status = http.StatusNoContent
			}
			return
		}
		// An error status with an untouched body falls back to the status text.
		if res.Status >= http.StatusBadRequest && reflect.ValueOf(&response.Body).Elem().IsZero() {
			return
		}
		res.Body = response.Body
	})
}

// Get registers a route handler for GET requests.
func (app *Application) Get[In, Out any](path string, handler RouteHandler[In, Out]) {
	app.Route("GET", path, handler)
}

// Head registers a route handler for HEAD requests.
func (app *Application) Head[In, Out any](path string, handler RouteHandler[In, Out]) {
	app.Route("HEAD", path, handler)
}

// Post registers a route handler for POST requests.
func (app *Application) Post[In, Out any](path string, handler RouteHandler[In, Out]) {
	app.Route("POST", path, handler)
}

// Put registers a route handler for PUT requests.
func (app *Application) Put[In, Out any](path string, handler RouteHandler[In, Out]) {
	app.Route("PUT", path, handler)
}

// Patch registers a route handler for PATCH requests.
func (app *Application) Patch[In, Out any](path string, handler RouteHandler[In, Out]) {
	app.Route("PATCH", path, handler)
}

// Delete registers a route handler for DELETE requests.
func (app *Application) Delete[In, Out any](path string, handler RouteHandler[In, Out]) {
	app.Route("DELETE", path, handler)
}

// Options registers a route handler for OPTIONS requests.
func (app *Application) Options[In, Out any](path string, handler RouteHandler[In, Out]) {
	app.Route("OPTIONS", path, handler)
}

// Trace registers a route handler for TRACE requests.
func (app *Application) Trace[In, Out any](path string, handler RouteHandler[In, Out]) {
	app.Route("TRACE", path, handler)
}
