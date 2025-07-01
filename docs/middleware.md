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

Vox's core concept is the middleware system. You can think of a vox application as a chain of middlewares. When a request comes in, the middlewares will be executed one by one.

The middleware can pre-process the request, for example, by extracting cookies from the HTTP header, transforming them into a user or session object, and storing the result in the context for future use.

A middleware can also terminate the execution of the next middlewares and respond to the user. This is useful for authentication or input validation.

Middleware can also modify the request or response. You can parse input data from JSON to a Go struct for a known schema, so you don't need to process it in your main business handler. You can also marshal the result/error to JSON or other encoding types in one place.

Your actual business handler can also be a middleware, and this is usually intended to be the last in the middleware chain.

## A basic middleware

The simplest middleware changes the response body to a string like this:

```go
func(ctx *vox.Context, req *vox.Request, res *vox.Response) {
    res.Body = "Hello, world!"
}
```

The `res.Body` will be written to the response HTTP body. If someone opens your website, they should see the string you wrote.

## Middleware for pre/post-processing

Here is an example of a middleware that records the current time, calls the next middleware, and then modifies the response to include the total processing time for the current request in the HTTP header.

Please notice the `ctx.Next()` call. This call moves the execution to the next middleware in the chain. When the next middleware finishes, `ctx.Next()` will return.

The `ctx.Next()` function takes no arguments and has no return value. Input and output should be handled through the `Request`, `Response`, and `Context` objects, or through global variables if you prefer.

```go
func(ctx *vox.Context, req *vox.Request, res *vox.Response) {
    start := time.Now()
    ctx.Next()
    duration := time.Now().Sub(start)
    res.Header.Set("X-Response-Time", fmt.Sprintf("%s", duration))
}
```

## Terminate execution

This is a simple validation example. You can validate a token in the request's HTTP header. If it's validated, call `ctx.Next()` to continue to the next middleware. Otherwise, set the status code and body in the response for an error message and return to the previous middlewares until the top to respond to the user.

```go
func(ctx *vox.Context, req *vox.Request, res *vox.Response) {
    if req.Header.Get("X-API-Token") == "a-secret" {
        ctx.Next()
    }
    res.Status = 403
    res.Body = "You shall not pass!"
}
```
