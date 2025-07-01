---
title: Context
nav_order: 4
---

# Context

Vox's `Context` object is a wrapper around the standard `context.Context` from the Go standard library. It provides a way to pass data between middlewares and to control the execution of the middleware chain.

## App

The `App` field is a pointer to the `vox.Application` instance. This can be used to access the application's configuration or other properties.

## Next

The `Next` function is used to call the next middleware in the chain. It's the middleware's responsibility to call the `Next` function. If a middleware does not call `Next`, the execution of the middleware chain will be terminated.

Here is an example of a simple logging middleware that uses the `Context` object to pass data to the next middleware:

```go
func Logger(ctx *vox.Context, req *vox.Request, res *vox.Response) {
    start := time.Now()
    ctx.Next()
    log.Printf("%s %s %v", req.Method, req.URL.Path, time.Since(start))
}
```
