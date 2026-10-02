---
title: Usage
nav_order: 100
---

# Usage
{: .no_toc }

## Table of contents
{: .no_toc .text-delta }

1. TOC
{:toc}

---

## Quick review

```go
package main

import (
	"fmt"
	"time"

	"github.com/aisk/vox"
)

func main() {
	app := vox.New()

	// custom middleware that adds an X-Response-Time header
	app.Use(func(ctx *vox.Context, req *vox.BaseRequest, res *vox.BaseResponse) {
		start := time.Now()
		ctx.Next()
		duration := time.Now().Sub(start)
		res.Header.Set("X-Response-Time", fmt.Sprintf("%s", duration))
	})

	// router param
	app.Get("/hello/{name}", func(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[string]) {
		res.Body = "Hello, " + req.Params["name"] + "!"
	})

	// response
	app.Get("/", func(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[string]) {
		// get the query string
		name := req.URL.Query().Get("name")
		if name == "" {
			name = "World"
		}
		res.Body = "Hello, " + name + "!"
	})

	app.Run("localhost:3000")
}
```

## Handle HTTP Methods

```go
package main

import (
	"github.com/aisk/vox"
)

func handler(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[string]) {
	// Get the current request's HTTP method and put it to the result page.
	res.Body = "HTTP Method is: " + req.Method
}

func main() {
	app := vox.New()

	app.Get("/", handler)
	app.Post("/", handler)
	app.Put("/", handler)
	app.Patch("/", handler)
	app.Delete("/", handler)
	app.Head("/", handler)
	app.Options("/", handler)
	app.Trace("/", handler)

	// In some cases you may need to handle a custom HTTP method not defined by RFCs, such as FLY.
	app.Route("FLY", "/", handler)

	app.Run("localhost:3000")
}
```

## Match any HTTP method

You can match all HTTP methods for a path with `"*"`:

```go
package main

import (
	"github.com/aisk/vox"
)

func main() {
	app := vox.New()
	app.Route("*", "/health", func(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[string]) {
		res.Body = "ok: " + req.Method
	})
	app.Run("localhost:3000")
}
```

Use this pattern for generic endpoints (for example, debug probes). For business APIs, explicit methods (`Get`, `Post`, `Put`, ...) are usually clearer.

## Serve static files

You can expose a local directory under a URL prefix with static middleware:

```go
package main

import (
	"github.com/aisk/vox"
	"github.com/aisk/vox/middlewares/static"
)

func main() {
	app := vox.New()

	// Serve files in ./public under /assets
	// GET /assets/logo.png -> ./public/logo.png
	app.Use(static.Middleware("/assets", "./public"))

	app.Run("localhost:3000")
}
```

## Get route parameters in URL path

```go
package main

import (
	"github.com/aisk/vox"
)

func hello(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[string]) {
	name := req.Params["name"]
	res.Body = "Hello, " + name + "!"
}

func main() {
	app := vox.New()
	app.Get("/hello/{name}", hello)
	app.Run("localhost:3000")
}
```

## Route matching

Routes use the pattern syntax of Go's `http.ServeMux`. When several patterns match a request, the most specific one wins, regardless of registration order, so `/users/me` takes precedence over `/users/{id}`.

A pattern ending in `/` matches every path under it. In particular `/` matches all paths, so once it is registered no request gets a 404. Use `/{$}` to match only the root path, and `{name...}` to capture the rest of the path:

```go
app.Get("/{$}", index)                // only "/"
app.Get("/files/{path...}", download) // "/files/a/b.txt" gives req.Params["path"] == "a/b.txt"
```

## Get query string parameters in URL

```go
package main

import (
	"github.com/aisk/vox"
)

func hello(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[string]) {
	name := req.URL.Query().Get("name")
	res.Body = "Hello, " + name + "!"
}

func main() {
	app := vox.New()
	app.Get("/hello", hello)
	app.Run("localhost:3000")
}
```

## Set response data

```go
package main

import (
	"github.com/aisk/vox"
)

func towel(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[string]) {
	// Set the response body.
	res.Body = "new towel is created!"
	// Set the response status code.
	res.Status = 201
	// Set the response header.
	res.Header.Set("Location", "/towels/42")
}

func main() {
	app := vox.New()
	app.Post("/towels", towel)
	app.Run("localhost:3000")
}
```

## Processing JSON requests and responses

Go 1.27 or later is required. Route type arguments are inferred from the handler.

```go
package main

import "github.com/aisk/vox"

type Towel struct {
    Color string `json:"color"`
    Size  string `json:"size"`
}

func towel(ctx *vox.Context, req *vox.Request[Towel], res *vox.Response[Towel]) {
    res.Body = req.Body
    res.Status = 201
    res.Header.Set("Location", "/towels/42")
}

func main() {
    app := vox.New()
    app.Post("/towels", towel)
    app.Run("localhost:3000")
}
```

Inputs other than `vox.NoBody` are decoded as one JSON value before the handler runs. Unsupported content types produce 415; bodies over the size limit produce 413; malformed, empty, `null`, or incompatible JSON and trailing data produce 400. JSON field validation is the application's responsibility. Use `Request[vox.NoBody]` for uploads or manual decoding, and read the original stream through `req.Request.Body`.

See [Request](request.md) for decoding behavior and [Response](response.md) for output formats and status codes.
