---
title: Error Handling
parent: Guide
nav_order: 7
---

# Error Handling
{: .no_toc }

Vox has no special error type or error return value. A failure is a response like any other: a status and, optionally, a body.

1. TOC
{:toc}

---

## Errors in route handlers

### Status only

Set a status of 400 or above and return. As long as `res.Body` still has its zero value, the client receives the status text:

```go
func showUser(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[User]) {
	user, ok := findUser(req.Params["id"])
	if !ok {
		res.Status = 404 // "Not Found"
		return
	}
	res.Body = user
}
```

### Status with a message or payload

The output type of a route describes its successful response. To send something else, assign it to `res.BaseResponse.Body`, which accepts any value:

```go
func showUser(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[User]) {
	user, err := loadUser(req.Params["id"])
	if err != nil {
		res.Status = 502
		res.BaseResponse.Body = map[string]string{"error": "user service unavailable"}
		return
	}
	res.Body = user
}
```

### Go errors

An `error` used as a body is written as its message, and the status defaults to 500:

```go
func showUser(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[User]) {
	user, err := loadUser(req.Params["id"])
	if err != nil {
		res.BaseResponse.Body = err // 500 with err.Error() as the body
		return
	}
	res.Body = user
}
```

{: .warning }
The message of an internal error can reveal details you do not want to expose. Prefer the middleware below, which logs the error and sends a generic message.

## Formatting errors in one place

Because middleware sees the response after `ctx.Next()`, error formatting can live in a single middleware. This one turns every error body into a JSON object, and hides the message of unexpected errors:

```go
func jsonErrors(ctx *vox.Context, req *vox.BaseRequest, res *vox.BaseResponse) {
	ctx.Next()

	err, ok := res.Body.(error)
	if !ok {
		return
	}
	if res.Status == 0 {
		res.Status = 500
	}
	message := err.Error()
	if res.Status >= 500 {
		log.Printf("%s %s: %v", req.Method, req.URL.Path, err)
		message = http.StatusText(res.Status)
	}
	res.Body = map[string]string{"error": message}
}
```

Handlers then report failures by assigning an error and, when 500 is not appropriate, a status. The check for `res.Status == 0` is needed because the default status is only applied after your middleware has returned.

## Request decoding errors

When a JSON request body is rejected, the route handler is not called. Vox responds with status 400, 413 or 415 and a short plain text message. The cases are listed in [Request]({% link guide/request.md %}#json-bodies).

Middleware sees the failure as a `*vox.DecodeError`. Since `*vox.DecodeError` implements `error` and carries its status, the `jsonErrors` middleware above already converts it. To handle it separately, check for the type:

```go
func decodeErrors(ctx *vox.Context, req *vox.BaseRequest, res *vox.BaseResponse) {
	ctx.Next()
	if err, ok := res.Body.(*vox.DecodeError); ok {
		res.Body = map[string]any{
			"error":  err.Message,
			"status": err.Status,
		}
	}
}
```

## Not found

A request that matches no route gets `404 Not Found` as plain text. To customize the response, add a middleware that fills it in when nothing else has:

```go
func notFound(ctx *vox.Context, req *vox.BaseRequest, res *vox.BaseResponse) {
	ctx.Next()
	if !res.HasBody() && res.Status == 0 {
		res.Status = 404
		res.Body = map[string]string{"error": "no such route: " + req.URL.Path}
	}
}
```

Both conditions matter. `res.HasBody()` is false when no route has responded, and `res.Status == 0` leaves alone the routes that deliberately set an error status without a body.

A request whose path is registered only for other methods is also treated as unmatched.

## Panics

Vox does not recover from panics. A panic in a handler propagates to Go's HTTP server, which logs the stack trace and closes the connection without sending a response. To send a 500 response, recover in a middleware:

```go
func recovery(ctx *vox.Context, req *vox.BaseRequest, res *vox.BaseResponse) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic: %v\n%s", r, debug.Stack())
			res.Status = 500
			res.Body = errors.New("internal server error")
		}
	}()
	ctx.Next()
}
```

It catches panics from the middleware registered after it and from route handlers, so register it early.

{: .note }
Writing the response happens outside of your middleware, in the built-in `respond` middleware. A failure at that stage, such as a body value that cannot be encoded as JSON, is not caught by a recovery middleware.

## Putting it together

Order matters when these are combined. Each middleware handles what the ones registered after it leave behind:

```go
app := vox.New()
app.Use(jsonErrors) // formats error bodies, including the one set by recovery
app.Use(recovery)   // turns panics into a 500 error
app.Use(notFound)   // answers requests that no route matched
```
