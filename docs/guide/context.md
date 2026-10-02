---
title: Context
parent: Guide
nav_order: 4
redirect_from:
  - /docs/context
---

# Context
{: .no_toc }

Every middleware and route handler receives a `*vox.Context`. One context is created per request and shared by the whole chain.

1. TOC
{:toc}

---

## Fields

```go
type Context struct {
	context.Context
	App  *Application
	Next func()
}
```

| Field | Purpose |
|:--|:--|
| `Context` | The embedded standard context, initially the one of the incoming `*http.Request`. |
| `App` | The application handling the request. |
| `Next` | Runs the rest of the middleware chain. For middleware only. |

## A standard context

`*vox.Context` embeds `context.Context`, so it can be passed to any function that expects one:

```go
func showUser(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[User]) {
	row := db.QueryRowContext(ctx, "SELECT id, name FROM users WHERE id = ?", req.Params["id"])
	// ...
}
```

The context is canceled when the client disconnects, so work that respects it stops early. Check `ctx.Done()` or `ctx.Err()` in long running handlers:

```go
select {
case result := <-work:
	res.Body = result
case <-ctx.Done():
	res.Status = 503
}
```

## Passing values

To make a value available to the rest of the chain, replace the embedded context with a derived one. All handlers of a request share the same `*vox.Context`, so the change is visible to every middleware and route handler that runs afterwards.

```go
type userKey struct{}

func authenticate(ctx *vox.Context, req *vox.BaseRequest, res *vox.BaseResponse) {
	user, ok := lookupUser(req.Header.Get("Authorization"))
	if !ok {
		res.Status = 401
		return
	}
	ctx.Context = context.WithValue(ctx.Context, userKey{}, user)
	ctx.Next()
}

func profile(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[User]) {
	res.Body = ctx.Value(userKey{}).(User)
}
```

As with any `context.WithValue` call, use an unexported key type to avoid collisions between packages.

{: .note }
Always read values from `ctx`. The context returned by `req.Context()` belongs to the original `*http.Request` and does not see values or deadlines added to `ctx.Context`.

## Deadlines

The same technique sets a deadline for everything that runs after a middleware:

```go
func timeout(d time.Duration) vox.Handler {
	return func(ctx *vox.Context, req *vox.BaseRequest, res *vox.BaseResponse) {
		var cancel context.CancelFunc
		ctx.Context, cancel = context.WithTimeout(ctx.Context, d)
		defer cancel()
		ctx.Next()
	}
}
```

The deadline does not interrupt a handler by itself. It takes effect in code that observes the context, such as database drivers and HTTP clients that were given `ctx`.

## App

`ctx.App` gives handlers access to the application, most commonly to read [configuration]({% link guide/configuration.md %}):

```go
func version(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[string]) {
	res.Body = ctx.App.GetConfig("app:version")
}
```

## Next

`ctx.Next()` runs the next middleware in the chain, and eventually the router. A middleware decides whether and when to call it. If it does not, the rest of the chain is skipped.

```go
func logger(ctx *vox.Context, req *vox.BaseRequest, res *vox.BaseResponse) {
	start := time.Now()
	ctx.Next()
	log.Printf("%s %s %v", req.Method, req.URL.Path, time.Since(start))
}
```

Call `ctx.Next()` at most once per middleware. `Next` is for middleware only: a route handler is the end of the chain and simply returns, and calling `Next` from one panics. See [Middleware]({% link guide/middleware.md %}) for the full picture.
