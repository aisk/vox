---
title: Middleware
parent: Guide
nav_order: 3
redirect_from:
  - /docs/middleware
---

# Middleware
{: .no_toc }

Middleware is the core concept of Vox. An application is a chain of middleware with the router at its end, and every request passes through the chain.

1. TOC
{:toc}

---

## Signature

A middleware is a `vox.Handler`:

```go
func(ctx *vox.Context, req *vox.BaseRequest, res *vox.BaseResponse)
```

Register it with `app.Use`. Unlike route handlers, middleware works with `BaseRequest` and `BaseResponse`, the untyped views of the request and response, because it runs for every route.

## How the chain works

Each middleware wraps everything registered after it. Calling `ctx.Next()` runs the rest of the chain and returns when it has finished:

```
logging
└─ respond
   └─ middleware A
      └─ middleware B
         └─ router
            └─ route handler
```

A middleware can do three things:

- Run code **before** `ctx.Next()`, for example to authenticate the request or to store a value in the context.
- Run code **after** `ctx.Next()`, for example to add a header, record metrics, or replace the response.
- **Not call** `ctx.Next()` at all. The rest of the chain, including the route handler, is skipped and the middleware's own response is sent.

`ctx.Next()` takes no arguments and returns nothing. Input and output travel through `ctx`, `req` and `res`.

## Execution order

`vox.New()` starts the chain with two built-in middleware, `logging` and `respond`. Middleware registered with `app.Use` follows in registration order. The router always runs last, no matter where routes are registered relative to `app.Use`.

```go
app := vox.New()
app.Use(middlewareA)
app.Get("/", handler)
app.Use(middlewareB)
```

A request for `/` flows like this:

1. `logging`
2. `respond`
3. `middlewareA`, up to its `ctx.Next()`
4. `middlewareB`, up to its `ctx.Next()`
5. the router and the matched route handler
6. the rest of `middlewareB`
7. the rest of `middlewareA`
8. `respond` writes the status, headers and body to the client
9. `logging` prints the access log line

## Writing middleware

### Before and after

This middleware records the time, runs the rest of the chain, and then adds the elapsed time to the response headers:

```go
func responseTime(ctx *vox.Context, req *vox.BaseRequest, res *vox.BaseResponse) {
	start := time.Now()
	ctx.Next()
	res.Header.Set("X-Response-Time", time.Since(start).String())
}
```

Setting a header after `ctx.Next()` works because nothing is sent to the client until the chain has unwound back to `respond`.

### Short-circuiting

This middleware rejects requests without a valid token. When the check fails it sets a response and returns without calling `ctx.Next()`, so no route runs:

```go
func requireToken(ctx *vox.Context, req *vox.BaseRequest, res *vox.BaseResponse) {
	if req.Header.Get("X-API-Token") != "a-secret" {
		res.Status = 403
		res.Body = "You shall not pass!"
		return
	}
	ctx.Next()
}
```

### Responding directly

A middleware that never calls `ctx.Next()` answers every request, and no route is ever reached:

```go
func maintenance(ctx *vox.Context, req *vox.BaseRequest, res *vox.BaseResponse) {
	res.Status = 503
	res.Body = "Down for maintenance"
}
```

### Options

A middleware that needs configuration is usually written as a function that returns a `vox.Handler`:

```go
func requireHeader(name, value string) vox.Handler {
	return func(ctx *vox.Context, req *vox.BaseRequest, res *vox.BaseResponse) {
		if req.Header.Get(name) != value {
			res.Status = 403
			return
		}
		ctx.Next()
	}
}

app.Use(requireHeader("X-API-Token", "a-secret"))
```

### Limiting middleware to some paths

`app.Use` applies to every request. To restrict a middleware, check the request inside it:

```go
func adminOnly(ctx *vox.Context, req *vox.BaseRequest, res *vox.BaseResponse) {
	if strings.HasPrefix(req.URL.Path, "/admin/") && !isAdmin(req) {
		res.Status = 403
		return
	}
	ctx.Next()
}
```

## What middleware can see

The request is only matched to a route when the chain reaches the router, so some information is not available until `ctx.Next()` has returned.

| | Before `ctx.Next()` | After `ctx.Next()` |
|:--|:--|:--|
| `req.Params` | empty | filled by the matched route |
| `res.HasBody()` | `false` | `true` when a route or an inner middleware responded |
| `res.Body` | unset | the response body, for example the route's `Out` value or a `*vox.DecodeError` |
| `res.Status` | `0` | the status set explicitly, or `0` when it was left to the default |
| `res.Header` | writable | writable |

A status of `0` means "not set yet". The default (200, 204, 404 or 500, see [Response]({% link guide/response.md %}#status)) is filled in by `respond` after your middleware has returned.

### Inspecting and replacing the response

After `ctx.Next()`, `res.Body` holds the value the route committed. Middleware can inspect it with a type assertion and replace it:

```go
func envelope(ctx *vox.Context, req *vox.BaseRequest, res *vox.BaseResponse) {
	ctx.Next()
	if user, ok := res.Body.(User); ok {
		res.Body = map[string]any{"data": user}
	}
}
```

### Acting as a fallback

When no route matches, the router leaves the response untouched. A middleware can detect this after `ctx.Next()` and provide the response itself:

```go
func notFound(ctx *vox.Context, req *vox.BaseRequest, res *vox.BaseResponse) {
	ctx.Next()
	if !res.HasBody() && res.Status == 0 {
		res.Status = 404
		res.Body = map[string]string{"error": "no such route"}
	}
}
```

## Sharing data with handlers

Middleware passes values to later middleware and to route handlers through the context. See [Context]({% link guide/context.md %}#passing-values).

## Built-in middleware

Two middleware are part of every application and always run first.

`logging` prints one line per request to standard output:

```
127.0.0.1:53412 - - [03/Oct/2026:10:00:00 +0000] "GET /hello/gopher HTTP/1.1" 200 14
```

The fields are the remote address, the user from the URL (normally `-`), the time, the request line, the status and the number of body bytes written. Turn it off with the `logging:disable` setting described in [Configuration]({% link guide/configuration.md %}).

`respond` writes the response to the client once the rest of the chain has returned. It applies the default status and content type, then writes the body according to its type. See [Response]({% link guide/response.md %}). A handler or middleware that writes to the connection itself sets `res.DontRespond` to make `respond` step aside.

Packages with optional middleware are listed under [Bundled Middleware]({% link middlewares/index.md %}).

## Using net/http middleware

An `Application` is an `http.Handler`, so middleware written for `net/http` can wrap the whole application:

```go
app := vox.New()
handler := http.TimeoutHandler(app, 10*time.Second, "timed out")
http.ListenAndServe("localhost:3000", handler)
```

Such a wrapper runs outside the Vox chain, before `logging`.

## Route handlers and Next

`ctx.Next()` is only for middleware. A route handler is the end of the chain. Returning from it hands control back to the middleware that wrapped it, and calling `ctx.Next()` from a route handler panics.
