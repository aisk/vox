---
title: Middleware
nav_order: 2
---

# Middleware
{: .no_toc }

## Table of contents
{: .no_toc .text-delta }

1. TOC
{:toc}

---

Vox's core concept is the middleware system. You can think of a Vox application as a chain of middleware. When a request comes in, each middleware is executed in order.

The middleware can pre-process the request, for example, by extracting cookies from the HTTP header, transforming them into a user or session object, and storing the result in the context for future use.

A middleware can also terminate execution of the remaining middleware and respond to the user. This is useful for authentication or input validation.

Middleware can also modify the request or response. You can parse input data from JSON to a Go struct for a known schema, so you don't need to process it in your main business handler. You can also marshal the result/error to JSON or other encoding types in one place.

Your actual business handler can also be a middleware, and this is usually intended to be the last in the middleware chain.

## Execution order

Vox has a built-in middleware chain. In `vox.New()`, the initial order is:

1. `logging`
2. `routeHandler`
3. `respond`

Middleware functions registered by `app.Use(...)` are appended after these built-ins.

For example:

```go
app := vox.New()
app.Use(middlewareA)
app.Use(middlewareB)
```

The execution flow is:

1. `logging`
2. `routeHandler`
3. `respond`
4. `middlewareA`
5. `middlewareB`
6. back to `respond` to write status/header/body to the client

Also, `routeHandler` always calls `ctx.Next()`, whether a route is matched or not. This means middleware functions added by `app.Use(...)` can still run as fallback handlers.

## A basic middleware

The simplest middleware changes the response body to a string like this:

```go
func(ctx *vox.Context, req *vox.Request, res *vox.Response) {
    res.Body = "Hello, world!"
}
```

The `res.Body` will be written to the response HTTP body. If someone opens your website, they should see the string you wrote.

## Middleware for pre/post-processing

Here is an example of a middleware that records the current time, calls the next middleware, and then modifies the response to include total processing time in the `X-Response-Time` header.

Note the `ctx.Next()` call. It moves execution to the next middleware in the chain. When the next middleware finishes, `ctx.Next()` returns.

The `ctx.Next()` function takes no arguments and has no return value. Input and output should be handled through the `Request`, `Response`, and `Context` objects.

```go
func(ctx *vox.Context, req *vox.Request, res *vox.Response) {
    start := time.Now()
    ctx.Next()
    duration := time.Now().Sub(start)
    res.Header.Set("X-Response-Time", fmt.Sprintf("%s", duration))
}
```

## Terminate execution

This is a simple validation example. Validate a token in the request header. If the token is valid, call `ctx.Next()` to continue. Otherwise, set an error status/body and return.

```go
func(ctx *vox.Context, req *vox.Request, res *vox.Response) {
    if req.Header.Get("X-API-Token") != "a-secret" {
        res.Status = 403
        res.Body = "You shall not pass!"
        return
    }
    ctx.Next()
}
```
