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

Vox's core concept is the middleware system. You can think of a vox application is a chain of middlewares, when a request came in, the middlewares will be executed one by one to the last.

The middleware can pre-process the request, for example, extract cookies from HTTP header, transform it to user object or session object, store the result in context for future usage.

And the middleware can terminate the execution of next middlewares and respond to users, for authentication or input validation scenarios.

Middleware can also modify the request or response. You can parse input data from JSON to go struct for known schema, for you do not need to process it in your main actual business handler. You can also marshall the result/error to JSON or other encoding types in one place.

Your actual business handler can be a middleware also, and this usually intends to be the last in the middleware chain.

## A basic middleware

The simplest middleware is change the response body to a string like this:

```go
func(ctx *vox.Context, req *vox.Request, res *vox.Response) {
    res.Body = "Hello, world!"
}
```

The `res.Body` should write to response HTTP body, if someone opened your website, you should see the string you wrote.

## Middleware for pre/post-process

There is an example for do something like record the current time, and call the next middlewares, and modify the response, wrote the whole processing time for current request to the HTTP header.

Please notice the `ctx.Next()`, in this call, the execution will be moved to next middlewares in chain, and when they finished, `ctx.Next()` will be returned.

The `ctx.Next()` takes no argument, and do have no return. or input or output should via the request/response/context or global variables if you like.

```go
func(ctx *vox.Context, req *vox.Request, res *vox.Response) {
    start := time.Now()
    ctx.Next()
    duration := time.Now().Sub(start)
    res.Header.Set("X-Response-Time", fmt.Sprintf("%s", duration))
}
```

## Terminate execution

This is a simple validation example, you can validate the token in the request HTTP header. If it's validated, call `ctx.Next` for future middleware executions. Otherwise, set the status code and body in response for an error message, and return to previous middlewares until the top to respond to users.

```go
func(ctx *vox.Context, req *vox.Request, res *vox.Response) {
    if req.Header.Get("X-API-Token") == "a-secret" {
        ctx.Next()
    }
    res.Status = 403
    res.Body = "You shall not pass!"
}
```
