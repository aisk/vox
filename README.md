# VOX

[![Go Reference](https://pkg.go.dev/badge/github.com/aisk/vox.svg)](https://pkg.go.dev/github.com/aisk/vox)
[![Build Status](https://travis-ci.org/aisk/vox.svg?branch=master)](https://travis-ci.org/aisk/vox)
[![Codecov](https://img.shields.io/codecov/c/github/aisk/vox.svg)](https://codecov.io/gh/aisk/vox)
[![Go Report Card](https://goreportcard.com/badge/github.com/aisk/vox)](https://goreportcard.com/report/github.com/aisk/vox)
[![Maintainability](https://api.codeclimate.com/v1/badges/d9a7d62ccc89b1752cf3/maintainability)](https://codeclimate.com/github/aisk/vox/maintainability)
[![Gitter chat](https://badges.gitter.im/go-vox/Lobby.png)](https://gitter.im/go-vox/Lobby)

A golang web framework for humans, inspired by [Koa](http://koajs.com) heavily.

![VoxLogo](https://cloudflare-ipfs.com/ipfs/QmUL4GF4HXhW6JUcNqVZBU1BwbJ2QULh81v5ZjZjPAWjnx)

## Getting started

### Installation

Using the `go get` power:

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

## More Docs

https://aisk.github.io/vox/

## Need Support?

If you need help for using vox, or have other questions, welcome to our [gitter chat room](https://gitter.im/go-vox/Lobby).

## About the Project

Vox is &copy; 2016-2020 by [aisk](https://github.com/aisk).

### License

Vox is distributed by a [MIT license](https://github.com/aisk/vox/tree/master/LICENSE).
