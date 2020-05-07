---
layout: default
title: Home
nav_order: 1
description: "A golang web framework for humans, inspired by Koa heavily."
permalink: /
---

# VOX: Go Web Framework for Humans
{: .fs-9 }

A golang web framework for humans, inspired by Koa heavily.
{: .fs-6 .fw-300 }

[Get started now](#getting-started){: .btn .btn-primary .fs-5 .mb-4 .mb-md-0 .mr-2 } [View it on GitHub](https://github.com/aisk/vox){: .btn .fs-5 .mb-4 .mb-md-0 }

---

## Getting started

### Installation

Using the `go get` power:

```sh
$ go get -u github.com/aisk/vox
```

### Basic Web Application

```sh
package main

import (
	"fmt"
	"time"

	"github.com/aisk/vox"
)

func main() {
	app := vox.New()

	// custom middleware that add a x-response-time to the response header
	app.Use(func(ctx *vox.Context, req *vox.Request, res *vox.Response) {
		start := time.Now()
		ctx.Next()
		duration := time.Now().Sub(start)
		res.Header.Set("X-Response-Time", fmt.Sprintf("%s", duration))
	})

	// router param
	app.Get("/hello/{name}", func(ctx *vox.Context, req *vox.Request, res *vox.Response) {
		res.Body = "Hello, " + req.Params["name"] + "!"
	})

	app.Run("localhost:3000")
}
```

## About the project

Vox is &copy; 2016-2020 by [aisk](https://github.com/aisk).

### License

Vox is distributed by a [MIT license](https://github.com/aisk/vox/tree/master/LICENSE).
