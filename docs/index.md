---
layout: home

hero:
  name: Vox
  text: A Go web framework for humans
  tagline: Koa style middleware and typed route handlers, on top of net/http.
  actions:
    - theme: brand
      text: Get Started
      link: /getting-started
    - theme: alt
      text: View on GitHub
      link: https://github.com/aisk/vox

features:
  - title: Typed route handlers
    details: The request body is decoded from JSON before your handler runs and the response body is encoded when it returns. Invalid input is rejected before it reaches your code.
  - title: A middleware chain
    details: Each middleware wraps everything after it. It can run code before and after the route handler, short-circuit the request, or rewrite the response.
  - title: Familiar routing
    details: Paths use the pattern syntax of Go's http.ServeMux, with named parameters, wildcards and most-specific-wins precedence.
  - title: No lock-in
    details: An Application is an http.Handler, so it works with http.Server, httptest and your existing net/http code.
---

## A quick look

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

Vox requires Go 1.27 or later.

## Where to go next

| If you want to | Read |
|:--|:--|
| Build and run a first application | [Getting Started](./getting-started.md) |
| Understand how requests are matched | [Routing](./guide/routing.md) |
| Understand typed handlers | [Route Handlers](./guide/handlers.md) |
| Write your own middleware | [Middleware](./guide/middleware.md) |
| Serve files or expose profiling endpoints | [Bundled Middleware](./middlewares/index.md) |
| Copy a working snippet for a common task | [Recipes](./recipes.md) |
| Look up a type or method | [API reference](https://pkg.go.dev/github.com/aisk/vox) |
