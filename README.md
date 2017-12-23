# VOX

[![GoDoc](https://godoc.org/github.com/aisk/vox?status.svg)](https://godoc.org/github.com/aisk/vox) [![Build Status](https://travis-ci.org/aisk/vox.svg?branch=master)](https://travis-ci.org/aisk/vox) [![Codecov](https://img.shields.io/codecov/c/github/aisk/vox.svg)](https://codecov.io/gh/aisk/vox)

A golang web framework for humans, inspired by [Koa](http://koajs.com) heavily.

![](https://i.v2ex.co/9MO3sMs4.jpeg)

---

## Example

```go
package main

import (
	"fmt"
	"time"
	"regexp"

	"github.com/aisk/vox"
)

func main() {
	app := vox.New()

	// x-response-time
	app.Use(func(ctx *vox.Context, next func()) {
		start := time.Now()
		next()
		ms := time.Now().Sub(start).Seconds() / 1000
		ctx.Response.Header.Set("X-Response-Time", fmt.Sprintf("%fms", ms))
	})

	// logger
	app.Use(func(ctx *vox.Context, next func()) {
		next()
		fmt.Printf("%s %s\n", ctx.Request.Method, ctx.Request.URL)
	})

	// response
	app.Get(func(ctx *vox.Context) {
		ctx.Response.SetBody("Hello, World!")
	})

	// router param
	app.Get(regexp.MustCompile(`/hello/(?P<name>\w+)`), func(ctx *vox.Context) {
		ctx.Response.SetBody("Hello, " + ctx.Request.Params["name"] + "!")
	})

	app.Run(":3000")
}
```
