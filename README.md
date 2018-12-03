# VOX

[![GoDoc](https://godoc.org/github.com/aisk/vox?status.svg)](https://godoc.org/github.com/aisk/vox)
[![Build Status](https://travis-ci.org/aisk/vox.svg?branch=master)](https://travis-ci.org/aisk/vox)
[![Codecov](https://img.shields.io/codecov/c/github/aisk/vox.svg)](https://codecov.io/gh/aisk/vox)
[![Go Report Card](https://goreportcard.com/badge/github.com/aisk/vox)](https://goreportcard.com/report/github.com/aisk/vox)
[![Maintainability](https://api.codeclimate.com/v1/badges/d9a7d62ccc89b1752cf3/maintainability)](https://codeclimate.com/github/aisk/vox/maintainability)
[![Gitter chat](https://badges.gitter.im/go-vox/Lobby.png)](https://gitter.im/go-vox/Lobby)

A golang web framework for humans, inspired by [Koa](http://koajs.com) heavily.

![](https://i.v2ex.co/9MO3sMs4.jpeg)

---

## Installation

Using the `go get` power:

```sh
$ go get -u github.com/aisk/vox
```

---

## Example

### Quick review

```go
package main

import (
	"fmt"
	"time"

	"github.com/aisk/vox"
	"github.com/aisk/vox/middlewares/logging"
)

func main() {
	app := vox.New()

	// logging
	app.Use(logging.Middleware)

	// custom middleware that add a x-response-time to the response header
	app.Use(func(req *vox.Request, res *vox.Response) {
		start := time.Now()
		req.Next()
		duration := time.Now().Sub(start)
		res.Header.Set("X-Response-Time", fmt.Sprintf("%s", duration))
	})

	// router param
	app.Get("/hello/{name}", func(req *vox.Request, res *vox.Response) {
		res.Body = "Hello, " + req.Params["name"] + "!"
	})

	// response
	app.Get("/", func(req *vox.Request, res *vox.Response) {
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

### Handle HTTP Methods

```go
package main

import (
	"github.com/aisk/vox"
)

func handler(req *vox.Request, res *vox.Response) {
	// Get the current request's HTTP method and put it to the result page.
	res.Body = "HTTP Method is: " + req.Method
}

func main() {
	app := vox.New()

	app.Get("/", handler)
	app.Post("/", handler)
	app.Put("/", handler)
	app.Delete("/", handler)
	app.Head("/", handler)
	app.Options("/", handler)
	app.Trace("/", handler)

	// In some case you need handle custom HTTP method that not in the RFCs like FLY.
	app.Route("FLY", "/", handler)

	app.Run("localhost:3000")
}
```


---

## License

MIT
