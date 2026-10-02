---
title: Context
nav_order: 4
---

# Context

Vox's `Context` object is a wrapper around the standard `context.Context` from the Go standard library. It provides a way to pass data between middleware and control the execution of the middleware chain.

## App

The `App` field is a pointer to the `vox.Application` instance. This can be used to access the application's configuration or other properties.

## Next

The `Next` function is used to call the next middleware in the chain. It's the middleware's responsibility to call the `Next` function. If a middleware does not call `Next`, the execution of the middleware chain will be terminated.

`Next` is for middleware only. A route handler is the end of the chain and simply returns; calling `Next` from one panics on the first request to that route.

Middleware receives `BaseRequest` and `BaseResponse`. This logger calls `ctx.Next()` and measures the time spent in subsequent middleware and the route handler:

```go
func Logger(ctx *vox.Context, req *vox.BaseRequest, res *vox.BaseResponse) {
    start := time.Now()
    ctx.Next()
    log.Printf("%s %s %v", req.Method, req.URL.Path, time.Since(start))
}
```
