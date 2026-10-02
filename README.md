# VOX

[![Go Reference](https://pkg.go.dev/badge/github.com/aisk/vox.svg)](https://pkg.go.dev/github.com/aisk/vox)
[![Build Status](https://github.com/aisk/vox/actions/workflows/go.yml/badge.svg?branch=master)](https://github.com/aisk/vox/actions/workflows/go.yml)
[![Codecov](https://img.shields.io/codecov/c/github/aisk/vox.svg)](https://codecov.io/gh/aisk/vox)
[![Go Report Card](https://goreportcard.com/badge/github.com/aisk/vox)](https://goreportcard.com/report/github.com/aisk/vox)
[![Maintainability](https://api.codeclimate.com/v1/badges/d9a7d62ccc89b1752cf3/maintainability)](https://codeclimate.com/github/aisk/vox/maintainability)
[![Gitter chat](https://badges.gitter.im/go-vox/Lobby.png)](https://gitter.im/go-vox/Lobby)

A Go web framework for humans, heavily inspired by [Koa](http://koajs.com).

![VoxLogo](https://cloudflare-ipfs.com/ipfs/QmUL4GF4HXhW6JUcNqVZBU1BwbJ2QULh81v5ZjZjPAWjnx)

## Getting started

### Installation

Requires **Go 1.27 or later**. Install with `go get`:

```sh
$ go get -u github.com/aisk/vox
```

### Basic Web Application

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

	app.Post("/users", func(ctx *vox.Context, req *vox.Request[CreateUser], res *vox.Response[User]) {
		res.Status = 201
		res.Body = User{ID: 1, Name: req.Body.Name}
	})

	app.Run("localhost:3000")
}
```

Route handlers declare their request and response body types in the function signature. Vox decodes request bodies as JSON; use `Request[vox.NoBody]` when no decoding is needed. Middleware receives `BaseRequest` and `BaseResponse`.

## Documentation

The full documentation lives at https://aisk.github.io/vox/.

- [Getting Started](https://aisk.github.io/vox/getting-started)
- [Routing](https://aisk.github.io/vox/guide/routing), [Route Handlers](https://aisk.github.io/vox/guide/handlers) and [Middleware](https://aisk.github.io/vox/guide/middleware)
- [Request](https://aisk.github.io/vox/guide/request), [Response](https://aisk.github.io/vox/guide/response) and [Error Handling](https://aisk.github.io/vox/guide/errors)
- [Static Files](https://aisk.github.io/vox/middlewares/static) and [Pprof](https://aisk.github.io/vox/middlewares/pprof) middleware
- [Recipes](https://aisk.github.io/vox/recipes)

## Need Support?

If you need help for using vox, or have other questions, welcome to our [gitter chat room](https://gitter.im/go-vox/Lobby).

## About the Project

Vox is &copy; 2016-2026 by [aisk](https://github.com/aisk).

### License

Vox is distributed by a [MIT license](https://github.com/aisk/vox/tree/master/LICENSE).
