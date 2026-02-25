---
title: Pprof Middleware
nav_order: 7
---

# Pprof Middleware
{: .no_toc }

## Table of contents
{: .no_toc .text-delta }

1. TOC
{:toc}

---

Vox provides an optional middleware package for Go's built-in pprof endpoints.

## Register middleware

```go
package main

import (
	"github.com/aisk/vox"
	"github.com/aisk/vox/middlewares/pprof"
)

func main() {
	app := vox.New()
	app.Use(pprof.Middleware)

	app.Get("/", func(ctx *vox.Context, req *vox.Request, res *vox.Response) {
		res.Body = "Hello, World!"
	})

	app.Run("localhost:3000")
}
```

## Endpoints

After registering the middleware, these pprof paths are available:

- `/debug/pprof`
- `/debug/pprof/cmdline`
- `/debug/pprof/profile`
- `/debug/pprof/symbol`
- `/debug/pprof/trace`

## How it works

For pprof paths, the middleware sets `res.DontRespond = true` and writes directly via the raw `http.ResponseWriter`.

## Production note

pprof endpoints can expose runtime and performance details. In production, enable them only for trusted networks, or protect them with authentication/authorization middleware.
