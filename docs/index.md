---
title: Home
nav_order: 1
description: "A Go web framework for humans, heavily inspired by Koa."
permalink: /
---

# Vox
{: .fs-9 }

A Go web framework for humans, heavily inspired by Koa.
{: .fs-6 .fw-300 }

[Get started]({% link getting-started.md %}){: .btn .btn-primary .fs-5 .mb-4 .mb-md-0 .mr-2 }
[View it on GitHub](https://github.com/aisk/vox){: .btn .fs-5 .mb-4 .mb-md-0 }

---

Vox is a small web framework built on two ideas. Requests flow through a chain of middleware, as in [Koa](https://koajs.com), and route handlers declare the types of their request and response bodies in their signature.

```go
package main

import (
	"fmt"
	"time"

	"github.com/aisk/vox"
)

type CreateUser struct {
	Name string `json:"name"`
}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	app := vox.New()

	// A middleware that adds an X-Response-Time header.
	app.Use(func(ctx *vox.Context, req *vox.BaseRequest, res *vox.BaseResponse) {
		start := time.Now()
		ctx.Next()
		res.Header.Set("X-Response-Time", fmt.Sprintf("%s", time.Since(start)))
	})

	// A route with a path parameter and a plain text response.
	app.Get("/hello/{name}", func(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[string]) {
		res.Body = "Hello, " + req.Params["name"] + "!"
	})

	// A route that decodes a JSON body and responds with JSON.
	app.Post("/users", func(ctx *vox.Context, req *vox.Request[CreateUser], res *vox.Response[User]) {
		res.Status = 201
		res.Body = User{ID: 1, Name: req.Body.Name}
	})

	app.Run("localhost:3000")
}
```

## What you get

- **Typed route handlers.** `Request[In]` is decoded from JSON before your handler runs, and `Response[Out]` is encoded when it returns. Invalid input is rejected with a proper status code before it reaches your code.
- **A middleware chain.** Each middleware wraps everything after it, so it can run code before and after the route handler, short-circuit the request, or rewrite the response.
- **Familiar routing.** Paths use the pattern syntax of Go's `http.ServeMux`, with `{name}` parameters, `{name...}` wildcards and most-specific-wins precedence.
- **No lock-in.** An `Application` is an `http.Handler`, so it works with `http.Server`, `httptest` and existing `net/http` code.

Vox requires Go 1.27 or later.

## Where to go next

| If you want to | Read |
|:--|:--|
| Build and run a first application | [Getting Started]({% link getting-started.md %}) |
| Understand how requests are matched | [Routing]({% link guide/routing.md %}) |
| Understand typed handlers | [Route Handlers]({% link guide/handlers.md %}) |
| Write your own middleware | [Middleware]({% link guide/middleware.md %}) |
| Serve files or expose profiling endpoints | [Bundled Middleware]({% link middlewares/index.md %}) |
| Copy a working snippet for a common task | [Recipes]({% link recipes.md %}) |
| Look up a type or method | [API reference](https://pkg.go.dev/github.com/aisk/vox) |

## License

Vox is &copy; 2016-2026 by [aisk](https://github.com/aisk) and distributed under the [MIT license](https://github.com/aisk/vox/blob/master/LICENSE).
